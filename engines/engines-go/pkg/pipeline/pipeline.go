package pipeline

import (
	"cognitivexr.at/cogstream/engines/pkg/engine"
	"context"
	"gocv.io/x/gocv"
)

// ============================== data holders

type Frame struct {
	FrameId int
	Mat     *gocv.Mat
}

type EngineResult struct {
	FrameId int
	Result  interface{}
}

// ============================= communication interfaces

type EngineResultWriter interface {
	WriteResult(result *EngineResult) error
}

type FrameWriter interface {
	WriteFrame(frame *Frame) error
}

type FramePacketWriter interface {
	WriteFramePacket(packet *engine.FramePacket) error
}

type FrameChannel chan *Frame
type FramePacketChannel chan *engine.FramePacket
type EngineResultChannel chan *EngineResult

func (ch FrameChannel) WriteFrame(f *Frame) error {
	ch <- f
	return nil
}

func (ch FramePacketChannel) WriteFramePacket(packet *engine.FramePacket) error {
	ch <- packet
	return nil
}

func (ch EngineResultChannel) WriteResult(r *EngineResult) error {
	ch <- r
	return nil
}

// ============================= function interfaces

type Scanner interface {
	Scan(context.Context) (*engine.FramePacket, error)
}

type Decoder interface {
	Decode(ctx context.Context, packet *engine.FramePacket, dest FrameWriter) error
}

type Transformer interface {
	Transform(ctx context.Context, src *Frame, dest FrameWriter) error
}

type Engine interface {
	Process(ctx context.Context, frame *Frame, writer EngineResultWriter) error
}

// ============================= pipeline

type Pipeline struct {
	Scanner     Scanner
	Decoder     Decoder
	Transformer Transformer
	Engine      Engine

	cancel context.CancelFunc
}

func (p *Pipeline) Cancel() {
	p.cancel()
}


// functional types for Pipeline interfaces

type ScannerFunction func(context.Context) (*engine.FramePacket, error)

func (s ScannerFunction) Scan(ctx context.Context) (*engine.FramePacket, error) {
	return s(ctx)
}

type DecoderFunction func(ctx context.Context, packet *engine.FramePacket) (*Frame, error)

func (d DecoderFunction) Decode(ctx context.Context, packet *engine.FramePacket, dest FrameWriter) error {
	frame, err := d(ctx, packet)
	if err != nil {
		return err
	}
	return dest.WriteFrame(frame)
}

type TransformerFunction func(ctx context.Context, src *Frame) (*Frame, error)

func (t TransformerFunction) Transform(ctx context.Context, src *Frame, dest FrameWriter) error {
	frame, err := t(ctx, src)
	if err != nil {
		return err
	}
	return dest.WriteFrame(frame)
}
