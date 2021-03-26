package decoder

import (
	"cognitivexr.at/cogstream/engines/pkg/engine"
	"cognitivexr.at/cogstream/engines/pkg/pipeline"
	"context"
	"gocv.io/x/gocv"
)

type Decoder func(packet *engine.FramePacket) (*gocv.Mat, error)

func (d Decoder) Decode(_ context.Context, packet *engine.FramePacket, dest pipeline.FrameWriter) error {
	mat, err := d(packet)
	if err != nil {
		return err
	}
	frame := &pipeline.Frame{FrameId: int(packet.Header.FrameId), Mat: mat}
	return dest.WriteFrame(frame)
}

var colorImageDecoder = func(packet *engine.FramePacket) (*gocv.Mat, error) {
	mat, err := gocv.IMDecode(packet.Data, gocv.IMReadColor)
	return &mat, err
}

func ColorImageDecoder() Decoder {
	return colorImageDecoder
}
