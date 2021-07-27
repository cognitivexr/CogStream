package main

import (
	"bufio"
	"cognitivexr.at/cogstream/pkg/webrtc"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os/exec"
	"strconv"
)

const (
	frameX      = 960
	frameY      = 720
	frameSize   = frameX * frameY * 3
	minimumArea = 3000
)

func main() {
	//create pipe for transforming video to suitable format
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

	ppl := webrtc.NewWebRtcPipeline(ffmpegIn, ffmpegOut, rdb, offerChannel, responseChannel)
	ppl.SetUpPeer()
	ppl.ListenForIceConnection()
	log.Printf("connected!")
	ppl.PrintDescription()
	ppl.RunSequential()

}
