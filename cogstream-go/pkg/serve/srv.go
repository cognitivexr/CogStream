package serve

import (
	"cognitivexr.at/cogstream/pkg/engine"
	"context"
	"log"
	"net"
	"sync"
)

func ServeEngineNetwork(ctx context.Context, network string, address string, factory engine.Factory) error {
	// TODO: close
	ln, err := net.Listen(network, address)
	defer ln.Close()

	if err != nil {
		return err
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

		log.Println("calling SequentialEngineHandler in new connection")
		go func() {
			err = SequentialEngineHandler(ctx, conn, factory)
			if err != nil {
				log.Printf("stream stopped: %v", err)
			}
		}()
	}

	return nil
}

func ServeEngineNetworkSingle(ctx context.Context, network string, address string, factory engine.Factory) error {
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

	return SequentialEngineHandler(ctx, conn, factory)
}
