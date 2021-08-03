package webrtc

import (
	"bufio"
	"cognitivexr.at/cogstream/pkg/engine"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os/exec"
	"strconv"
)

func ServeEngineNetwork(ctx context.Context, factory engine.Factory) error {

	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-pix_fmt", "bgr24", "-s", strconv.Itoa(frameX)+"x"+strconv.Itoa(frameY), "-f", "rawvideo", "pipe:1") //nolint
	ffmpegIn, _ := ffmpeg.StdinPipe()
	ffmpegOut, _ := ffmpeg.StdoutPipe()
	ffmpegErr, _ := ffmpeg.StderrPipe()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// kill off program if ffmpeg fails
	if err := ffmpeg.Start(); err != nil {
		panic(err)
	}

	// relay all ffmpeg errors to go program
	go func() {
		scanner := bufio.NewScanner(ffmpegErr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	offerChannel := "offer"
	responseChannel := "response"

	ppl := NewWebRtcPipeline(FfmpegConfiguration{
		In:  ffmpegIn,
		Out: ffmpegOut,
	}, &RedisConfiguration{Client: rdb, OfferC: offerChannel, AnswerC: responseChannel}, factory.NewEngine())
	ppl.SetUpPeer()
	ppl.ListenForIceConnection()
	ppl.PrintDescription()
	ppl.RunSequential()

	return nil
}
