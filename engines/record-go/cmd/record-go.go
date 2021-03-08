package main

import (
	"cognitivexr.at/cogstream/engines/record-go/pkg/stream"
	"gocv.io/x/gocv"
	"log"
	"net"
)

func Serve(network string, address string) {
	ln, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("accept connection on address", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("calling connection handler")

		ctx := stream.NewStreamContext()
		saveVideo(ctx, conn)
	}
}


func saveVideo(ctx stream.StreamContext, conn net.Conn) {
	frames := make(chan []byte, 30)
	images := make(chan gocv.Mat)

	go stream.ImageDecoder(ctx, frames, images)
	go stream.SaveVideo(ctx, images)
	stream.ConnectionHandler(ctx, conn, frames)

	close(frames)
	close(images)
}

func windowDisplay(ctx stream.StreamContext, conn net.Conn) {
	frames := make(chan []byte, 30)
	images := make(chan gocv.Mat)

	go stream.ImageDecoder(ctx, frames, images)
	go stream.ConnectionHandler(ctx, conn, frames)
	stream.WindowDisplay(ctx, images)

	close(frames)
	close(images)
}

func main() {
	Serve("tcp", "0.0.0.0:53210")
}