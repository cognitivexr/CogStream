package pipeline

import (
	"cognitivexr.at/cogstream/engines/pkg/engine"
	"context"
	"fmt"
)

type scannerPipe struct {
	Scanner
	ctx context.Context
	dst FramePacketWriter
}

func (s *scannerPipe) Next() error {
	packet, err := s.Scan(s.ctx)
	if err != nil {
		return err
	}

	return s.dst.WriteFramePacket(packet)
}

type decoderPipe struct {
	Decoder
	ctx context.Context
	dst FrameWriter
}

func (d *decoderPipe) WriteFramePacket(packet *engine.FramePacket) error {
	return d.Decode(d.ctx, packet, d.dst)
}

type transformerPipe struct {
	Transformer
	ctx context.Context
	dst FrameWriter
}

func (t *transformerPipe) WriteFrame(frame *Frame) error {
	return t.Transform(t.ctx, frame, t.dst)
}

type enginePipe struct {
	Engine
	ctx context.Context
	dst EngineResultWriter
}

func (e *enginePipe) WriteFrame(frame *Frame) error {
	return e.Process(e.ctx, frame, e.dst)
}

func (p *Pipeline) RunSequential(ctx context.Context, sink EngineResultWriter) error {
	if p.Scanner == nil {
		return fmt.Errorf("pipeline scanner is nil")
	}
	if p.Decoder == nil {
		return fmt.Errorf("pipeline decoder is nil")
	}
	if p.Transformer == nil {
		return fmt.Errorf("pipeline transformer is nil")
	}
	if p.Engine == nil {
		return fmt.Errorf("pipeline engine is nil")
	}

	ctx, cancelFunc := context.WithCancel(ctx)
	p.cancel = cancelFunc

	eng := &enginePipe{p.Engine, ctx, sink}
	transformer := &transformerPipe{p.Transformer, ctx, eng}
	decoder := &decoderPipe{p.Decoder, ctx, transformer}
	scanner := &scannerPipe{p.Scanner, ctx, decoder}

	for {
		err := scanner.Next()
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
