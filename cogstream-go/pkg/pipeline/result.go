package pipeline

import (
	"cognitivexr.at/cogstream/pkg/stream"
	"encoding/json"
	"fmt"
)

type EngineResultSerializer interface {
	Serialize(packet *stream.ResultPacket, result *EngineResult) error
}

type JsonResultSerializer struct{}

func (j *JsonResultSerializer) Serialize(packet *stream.ResultPacket, result *EngineResult) error {
	bytes, err := json.Marshal(result.Result)
	if err != nil {
		return fmt.Errorf("could not marshal result payload to json: %v", err)
	}
	packet.Data = bytes
	return nil
}

type SerializingResultWriter struct {
	w stream.ResultPacketWriter
	s EngineResultSerializer
}

func NewJsonResultWriter(w stream.ResultPacketWriter) EngineResultWriter {
	return &SerializingResultWriter{
		w: w,
		s: &JsonResultSerializer{},
	}
}

func (j *SerializingResultWriter) WriteResult(result *EngineResult) error {
	packet := new(stream.ResultPacket)

	seconds := result.Timestamp.Unix()
	nanos := result.Timestamp.UnixNano() - seconds

	packet.Header = stream.ResultPacketHeader{
		StreamId:         0,
		FrameId:          uint32(result.FrameId),
		TimestampSeconds: uint32(seconds),
		TimestampNanos:   uint32(nanos),
	}

	if err := j.s.Serialize(packet, result); err != nil {
		return fmt.Errorf("could not serialize result packet: %v", err)
	}

	packet.Header.DataLen = uint32(len(packet.Data))

	return j.w.WritePacket(packet)
}
