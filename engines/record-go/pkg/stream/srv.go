package stream

import (
	"gocv.io/x/gocv"
	"log"
	"net"
)

func SaveVideoHandler(conn net.Conn) {
	frames := make(chan []byte, 30)
	images := make(chan gocv.Mat)
	ctx := NewStreamContext() // FIXME: initialize stream context

	go ImageDecoder(ctx, frames, images)
	go SaveVideoSink(ctx, images)
	ConnectionHandler(ctx, conn, frames)

	close(frames)
	close(images)
}

func WindowDisplayHandler(conn net.Conn) {
	frames := make(chan []byte, 30)
	images := make(chan gocv.Mat)
	ctx := NewStreamContext() // FIXME: initialize stream context

	go ImageDecoder(ctx, frames, images)
	go ConnectionHandler(ctx, conn, frames)
	WindowDisplaySink(ctx, images)

	close(frames)
	close(images)
}

// ServeSingle creates a server socket, accepts one connection, and then closes the server socket before initializing
// the stream.
func ServeSingle(network string, address string) {
	ln, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("accept connection on address", address)

	conn, err := ln.Accept()
	if err != nil {
		log.Fatal(err)
	}
	ln.Close()
	log.Println("calling connection handler")

	SaveVideoHandler(conn)
}

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

		SaveVideoHandler(conn)
	}
}