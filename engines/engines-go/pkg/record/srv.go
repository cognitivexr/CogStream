package record

import (
	"cognitivexr.at/cogstream/engines/pkg/engine"
	"context"
	"gocv.io/x/gocv"
	"log"
	"net"
	"sync"
)

func SaveVideoHandler(conn net.Conn) error {
	frames := make(chan []byte, 30)
	images := make(chan gocv.Mat)

	ctx, err := engine.InitStreamContext(conn)
	if err != nil {
		log.Println("error initializing stream context", err)
		conn.Close()
		return err
	}

	go engine.ImageDecoder(ctx, frames, images)
	go SaveVideoSink(ctx, images)
	engine.ConnectionHandler(ctx, conn, frames)

	close(frames)
	close(images)
	return nil
}

// ServeSingle creates a server socket, accepts one connection, and then closes the server socket before initializing
// the engine.
func ServeSingle(ctx context.Context, network string, address string) error {
	ln, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	if started, ok := ctx.Value("started").(*sync.WaitGroup); started != nil && ok {
		log.Println("started engine serving")
		started.Done()
	}

	log.Println("accept connection on address", address)

	conn, err := ln.Accept()
	if err != nil {
		ln.Close()
		return err
	}
	ln.Close()
	log.Println("calling connection handler")

	return SaveVideoHandler(conn)
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
