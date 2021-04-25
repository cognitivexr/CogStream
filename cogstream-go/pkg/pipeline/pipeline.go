package pipeline

import (
	"cognitivexr.at/cogstream/pkg/stream"
	"context"
	"errors"
	"gocv.io/x/gocv"
	"io"
	"time"
)

// ============================== data holders

var Stop = errors.New("stop pipeline")

type Frame struct {
	FrameId int
	Mat     *gocv.Mat
}

type EngineResult struct {
	FrameId   int
	Timestamp time.Time
	Result    interface{}
}

// ============================= communication interfaces

type EngineResultWriter interface {
	WriteResult(result *EngineResult) error
}

type FrameWriter interface {
	WriteFrame(frame *Frame) error
}

type FramePacketWriter interface {
	WriteFramePacket(packet *stream.FramePacket) error
}

type FrameChannel chan *Frame
type FramePacketChannel chan *stream.FramePacket
type EngineResultChannel chan *EngineResult

func (ch FrameChannel) WriteFrame(f *Frame) error {
	ch <- f
	return nil
}

func (ch FramePacketChannel) WriteFramePacket(packet *stream.FramePacket) error {
	ch <- packet
	return nil
}

func (ch EngineResultChannel) WriteResult(r *EngineResult) error {
	ch <- r
	return nil
}

type noopEngineResultWriter struct{}

func (n *noopEngineResultWriter) WriteResult(result *EngineResult) error {
	return nil
}

var NoopEngineResultWriter EngineResultWriter = &noopEngineResultWriter{}

// ============================= function interfaces

type Scanner interface {
	Scan(context.Context) (*stream.FramePacket, error)
}

type Decoder interface {
	Decode(ctx context.Context, packet *stream.FramePacket, dest FrameWriter) error
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
	Results     EngineResultWriter

	cancel context.CancelFunc
}

func (p *Pipeline) ConfigureForStream(stream *stream.Stream) {
	stream.AcceptConfigurators(p.Scanner, p.Decoder, p.Transformer, p.Engine, p.Results)
}

func (p *Pipeline) Cancel() {
	p.cancel()
}

func (p *Pipeline) Close() error {
	err := tryClose(p.Engine)
	return err
}

// functional types for Pipeline interfaces

type ScannerFunction func(context.Context) (*stream.FramePacket, error)

func (s ScannerFunction) Scan(ctx context.Context) (*stream.FramePacket, error) {
	return s(ctx)
}

type DecoderFunction func(ctx context.Context, packet *stream.FramePacket) (*Frame, error)

func (d DecoderFunction) Decode(ctx context.Context, packet *stream.FramePacket, dest FrameWriter) error {
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

func tryClose(obj interface{}) error {
	if closer, ok := obj.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
