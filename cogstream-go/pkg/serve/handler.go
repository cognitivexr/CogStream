package serve

import (
	"cognitivexr.at/cogstream/pkg/decoder"
	"cognitivexr.at/cogstream/pkg/engine"
	"cognitivexr.at/cogstream/pkg/pipeline"
	"cognitivexr.at/cogstream/pkg/stream"
	"cognitivexr.at/cogstream/pkg/transform"
	"context"
	"fmt"
	"net"
)

type engineResultPrinter struct{}

func (e engineResultPrinter) WriteResult(_ *pipeline.EngineResult) error {
	fmt.Println("engine result received")
	return nil
}

func SequentialEngineHandler(ctx context.Context, conn net.Conn, factory engine.Factory) error {
	defer conn.Close()

	s, err := stream.NewStream(ctx, conn)
	if err != nil {
		return fmt.Errorf("error initializing stream: %v", err)
	}
	s.Metadata().EngineFormat = factory.Descriptor().Specification.InputFormat

	transformer, err := transform.BuildTransformer(s.Metadata().ClientFormat, s.Metadata().EngineFormat)
	if err != nil {
		return fmt.Errorf("error initializing transformation pipeline: %v", err)
	}

	p := pipeline.Pipeline{
		Scanner:     stream.NewFramePacketScanner(stream.NewFramePacketReader(s)),
		Decoder:     decoder.ColorImageDecoder(),
		Transformer: transformer,
		Engine:      factory.NewEngine(),
	}

	s.AcceptConfigurators(p)

	err = p.RunSequential(ctx, engineResultPrinter{}) // FIXME: how to know about how to return results?

	return err
}
