package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"cognitivexr.at/cogstream-webrtc-client/gst"
	"cognitivexr.at/cogstream-webrtc-client/util"
	"github.com/go-redis/redis/v8"
	"github.com/pion/webrtc/v3"
)

func waitForRemoteSessionDescription(rdb *redis.Client, channel string) string {
	ctx := context.Background()
	pubsub := rdb.Subscribe(ctx, channel)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		return strings.TrimSpace(msg.Payload)
	}
}

func main() {
	videoSrc := flag.String("video-src", "v4l2src device=/dev/video0 ! video/x-raw, width=1280, height=720 ! videoconvert ! queue", "GStreamer video src")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	offerChannel := "offer"
	responseChannel := "response"

	// Everything below is the pion-WebRTC API! Thanks for using it ❤️.

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	// Create a video track
	vp8Track, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion")
	if err != nil {
		panic(err)
	} else if _, err = peerConnection.AddTrack(vp8Track); err != nil {
		panic(err)
	}

	// Create an offer to send to the browser
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the offer in base64 so we can paste it in browser
	fmt.Println(util.Encode(*peerConnection.LocalDescription()))
	rdb.Publish(context.TODO(), offerChannel, util.Encode(*peerConnection.LocalDescription()))

	// Wait for the answer to be submitted via HTTP
	// TODO wait for answer from redis
	answer := webrtc.SessionDescription{}
	util.Decode(waitForRemoteSessionDescription(rdb, responseChannel), &answer)

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(answer)
	if err != nil {
		panic(err)
	}

	// Start pushing buffers on these tracks
	gst.CreatePipeline("vp8", []*webrtc.TrackLocalStaticSample{vp8Track}, *videoSrc).Start()

	// Block forever
	select {}
}
