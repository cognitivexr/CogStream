package record

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/engines/pkg/engine"
	"context"
	"gocv.io/x/gocv"
	"log"
	"net"
	"sync"
)

func SaveVideoHandler(ctx context.Context, e *engines.EngineDescriptor, conn net.Conn) error {
	frames := make(chan *engine.FramePacket, 30)
	images := make(chan gocv.Mat)

	metadata := engine.NewStreamMetadata()
	ctx = engine.NewStreamContext(ctx, metadata)

	metadata.EngineFormat = e.Specification.InputFormat

	if attrs, ok := engine.GetAttributes(ctx); ok {
		// this is a mess. target format is encoded in OperationSpec sent by the client
		if format, err := messages.FormatFromAttributes(attrs); err == nil {
			log.Printf("deserialized engine format from context attributes: %v", format)
			metadata.EngineFormat = format
		} else {
			log.Printf("error while reading attributes: %v", err)
		}
	}

	err := engine.InitStream(ctx, conn)
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
func ServeSingle(ctx context.Context, engine *engines.EngineDescriptor, network string, address string) error {
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

	return SaveVideoHandler(ctx, engine, conn)
}

func Serve(ctx context.Context, engine *engines.EngineDescriptor, network string, address string) {
	ln, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}

	if started, ok := ctx.Value("started").(*sync.WaitGroup); started != nil && ok {
		log.Println("started engine serving")
		started.Done()
	}

	log.Println("accept connection on address", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("calling connection handler")
		SaveVideoHandler(ctx, engine, conn)
	}
}
