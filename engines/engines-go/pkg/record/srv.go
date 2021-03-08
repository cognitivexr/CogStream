package record

import (
	"cognitivexr.at/cogstream/engines/pkg/engine"
	"gocv.io/x/gocv"
	"log"
	"net"
)

func SaveVideoHandler(conn net.Conn) {
	frames := make(chan []byte, 30)
	images := make(chan gocv.Mat)
	ctx := engine.NewStreamContext() // FIXME: initialize engine context

	go engine.ImageDecoder(ctx, frames, images)
	go SaveVideoSink(ctx, images)
	engine.ConnectionHandler(ctx, conn, frames)

	close(frames)
	close(images)
}

// ServeSingle creates a server socket, accepts one connection, and then closes the server socket before initializing
// the engine.
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