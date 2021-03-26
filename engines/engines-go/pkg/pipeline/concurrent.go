package pipeline

import (
	"context"
	"fmt"
)

func (p *Pipeline) RunConcurrent(ctx context.Context) error {
	return runConcurrentPipeline(ctx, p)
}

func runConcurrentPipeline(ctx context.Context, p *Pipeline) error {
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

	errs := make(chan error)
	packets := make(FramePacketChannel)
	decoded := make(FrameChannel)
	transformed := make(FrameChannel)
	results := make(EngineResultChannel)

	defer func() {
		close(errs)
		close(packets)
		close(decoded)
		close(transformed)
		close(results)
	}()

	go RunEngine(ctx, p.Engine, transformed, results, errs)
	go RunTransformer(ctx, p.Transformer, decoded, transformed, errs)
	go RunDecoder(ctx, p.Decoder, packets, decoded, errs)
	go RunScanner(ctx, p.Scanner, packets, errs)

	for {
		select {
		case e := <-errs:
			// TODO: fault tolerance
			return e
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// functions to run pipeline components as go functions that communicate via channels

func RunScanner(ctx context.Context, scanner Scanner, packets FramePacketWriter, errs chan<- error) {
	for {
		packet, err := scanner.Scan(ctx)
		if err != nil {
			errs <- err
		}

		err = packets.WriteFramePacket(packet)
		if err != nil {
			errs <- err
		}

		if ctx.Err() != nil {
			return
		}
	}
}

func RunDecoder(ctx context.Context, d Decoder, src FramePacketChannel, dst FrameWriter, errs chan<- error) {
	for {
		select {
		case framePacket := <-src:
			err := d.Decode(ctx, framePacket, dst)
			if err != nil {
				errs <- err
			}
		case <-ctx.Done():
			return
		}
	}
}

func RunTransformer(ctx context.Context, t Transformer, src FrameChannel, dst FrameWriter, errs chan<- error) {
	for {
		select {
		case frame := <-src:
			err := t.Transform(ctx, frame, dst)
			if err != nil {
				errs <- err
			}
		case <-ctx.Done():
			return
		}
	}
}

func RunEngine(ctx context.Context, e Engine, src FrameChannel, dst EngineResultWriter, errs chan<- error) {
	for {
		select {
		case frame := <-src:
			err := e.Process(ctx, frame, dst)
			if err != nil {
				errs <- err
			}
		case <-ctx.Done():
			return
		}
	}
}
