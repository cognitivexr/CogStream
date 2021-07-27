package webrtc

import (
	"cognitivexr.at/cogstream/pkg/pipeline"
	"cognitivexr.at/cogstream/pkg/util"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/ivfwriter"
	"gocv.io/x/gocv"
	"io"
	"log"
	"strings"
	"time"
)

const (
	frameX    = 960
	frameY    = 720
	frameSize = frameX * frameY * 3
)

type Pipeline struct {
	config          *webrtc.Configuration
	conn            *webrtc.PeerConnection
	ivf             *ivfwriter.IVFWriter
	out             io.Reader
	engine          pipeline.Engine
	rdb             *redis.Client
	offerChannel    string
	responseChannel string
}

func NewWebRtcPipeline(in io.Writer, out io.Reader, rdb *redis.Client, offerChannel string, responseChannel string) *Pipeline {
	config := &webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	ivfWriter, err := ivfwriter.NewWith(in)
	if err != nil {
		panic(err)
	}

	p := Pipeline{
		config:          config,
		conn:            nil,
		out:             out,
		ivf:             ivfWriter,
		rdb:             rdb,
		offerChannel:    offerChannel,
		responseChannel: responseChannel,
	}
	return &p
}

func (p *Pipeline) SetUpPeer() {
	var err error
	p.conn, err = webrtc.NewPeerConnection(*p.config)
	if err != nil {
		panic(err)
	}
	p.conn.OnTrack(p.writeToIvf)
}

func (p *Pipeline) writeToIvf(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
	// Send Keyframe every 3 seconds
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for range ticker.C {
			errSend := p.conn.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}})
			if errSend != nil {
				fmt.Println(errSend)
			}
		}
	}()
	log.Printf("track with type %d", track.PayloadType())

	// write RTP packets into video stream
	for {
		rtp, _, readErr := track.ReadRTP()
		if readErr != nil {
			panic(readErr)
		}

		if ivfErr := p.ivf.WriteRTP(rtp); ivfErr != nil {
			panic(ivfErr)
		}
	}
}

func (p *Pipeline) waitForRemoteSessionDescription() string {
	ctx := context.Background()
	pubsub := p.rdb.Subscribe(ctx, p.offerChannel)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		return strings.TrimSpace(msg.Payload)
	}
}

func (p *Pipeline) ListenForIceConnection() {
	p.conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("connection state change: %s", connectionState.String())
	})

	offer := webrtc.SessionDescription{}
	util.Decode(p.waitForRemoteSessionDescription(), &offer)

	err := p.conn.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	answer, err := p.conn.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(p.conn)

	err = p.conn.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	<-gatherComplete
}

func (p *Pipeline) PrintDescription() {
	fmt.Println(util.Encode(*p.conn.LocalDescription()))
	p.rdb.Publish(context.TODO(), p.responseChannel, util.Encode(*p.conn.LocalDescription()))
}

func (p *Pipeline) RunSequential() {
	window := gocv.NewWindow("received video")
	defer window.Close()

	for {
		buf := make([]byte, frameSize)
		if _, err := io.ReadFull(p.out, buf); err != nil {
			log.Printf("could not read frame: %v", err)
			continue
		}
		img, _ := gocv.NewMatFromBytes(frameY, frameX, gocv.MatTypeCV8UC3, buf)
		if img.Empty() {
			log.Printf("Empty Mat")
			continue
		}

		// TODO: plant in Engine here
		window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}
	}
}
