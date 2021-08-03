package pipeline

import (
	"cognitivexr.at/cogstream/pkg/stream"
	"log"
	"testing"
	"time"
)

type resultPacketChannel chan *stream.ResultPacket

func (r resultPacketChannel) WritePacket(packet *stream.ResultPacket) error {
	r <- packet
	return nil
}

func TestJsonResultWriter_WriteResult(t *testing.T) {
	ch := make(resultPacketChannel, 1)

	writer := NewJsonResultWriter(ch)

	dict := make(map[string]string)
	dict["foo"] = "bar"
	dict["baz"] = "ed"

	result := &EngineResult{
		FrameId:   42,
		Timestamp: time.Now(),
		Result:    dict,
	}

	err := writer.WriteResult(result)
	if err != nil {
		t.Errorf("unexpected error while writing result: %v", err)
	}

	packet := <-ch

	log.Printf("FrameId: %d", packet.Header.FrameId)
	if packet.Header.FrameId != 42 {
		t.Error("unexpected frame id")
	}

	log.Printf("JSON: %s", packet.Data)
	if string(packet.Data) != "{\"baz\":\"ed\",\"foo\":\"bar\"}" {
		t.Errorf("unexpected packet data %s", packet.Data)
	}

	close(ch)
}
