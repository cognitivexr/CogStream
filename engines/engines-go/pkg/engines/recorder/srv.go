package recorder

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/pkg/decoder"
	"cognitivexr.at/cogstream/pkg/pipeline"
	"cognitivexr.at/cogstream/pkg/stream"
	"cognitivexr.at/cogstream/pkg/transform"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
)

type engineResultPrinter struct{}

func (e engineResultPrinter) WriteResult(_ *pipeline.EngineResult) error {
	fmt.Println("engine result received")
	return nil
}

func SaveVideoHandler(ctx context.Context, e *engines.EngineDescriptor, conn net.Conn) error {
	defer conn.Close()

	s, err := stream.NewStream(ctx, conn)
	if err != nil {
		return fmt.Errorf("error initializing stream: %v", err)
	}
	s.Metadata().EngineFormat = e.Specification.InputFormat

	transformer, err := transform.BuildTransformer(s.Metadata().ClientFormat, s.Metadata().EngineFormat)
	if err != nil {
		return fmt.Errorf("error initializing transformation pipeline: %v", err)
	}

	p := pipeline.Pipeline{
		Scanner:     stream.NewFramePacketScanner(stream.NewFramePacketReader(s)),
		Decoder:     decoder.ColorImageDecoder(),
		Transformer: transformer,
		Engine:      NewEngine(),
	}

	err = p.RunSequential(ctx, engineResultPrinter{})

	return err
}

func Serve(ctx context.Context, descriptor *engines.EngineDescriptor, network string, address string) {
	ln, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}

	if started, ok := ctx.Value("started").(*sync.WaitGroup); started != nil && ok {
		log.Println("started descriptor serving")
		started.Done()
	}

	log.Println("accept connection on address", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("calling connection handler")

		err = SaveVideoHandler(ctx, descriptor, conn)
		if err != nil {
			log.Printf("stream stopped: %v", err)
		}
	}
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
