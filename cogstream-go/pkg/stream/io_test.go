package stream

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestResultPacketWriter_WritePacket(t *testing.T) {
	packet := new(ResultPacket)
	packet.Header = ResultPacketHeader{
		StreamId:         1,
		FrameId:          2,
		TimestampSeconds: 1619190770,
		TimestampNanos:   395916981,
		DataLen:          19,
	}

	packet.Data = []byte("{'value': 'foobar'}")

	buf := bytes.NewBuffer(make([]byte, 64))
	buf.Reset()

	writer := NewResultPacketWriter(buf)
	err := writer.WritePacket(packet)

	if err != nil {
		t.Error(err)
	}

	if buf.Len() != (ResultPacketHeaderSize + int(packet.Header.DataLen)) {
		t.Error("did not write all bytes into packet")
	}

	if 1 != binary.LittleEndian.Uint32(buf.Next(4)) {
		t.Error("expected stream ID to be 1")
	}
	if 2 != binary.LittleEndian.Uint32(buf.Next(4)) {
		t.Error("expected frameId to be 2")
	}
	if 1619190770 != binary.LittleEndian.Uint32(buf.Next(4)) {
		t.Error("unexpected value for TimestampSeconds")
	}
	if 395916981 != binary.LittleEndian.Uint32(buf.Next(4)) {
		t.Error("unexpected value for TimestampNanos")
	}
	if packet.Header.DataLen != binary.LittleEndian.Uint32(buf.Next(4)) {
		t.Error("unexpected value for DataLen")
	}

	data := buf.Bytes()
	json := string(data)
	if json != "{'value': 'foobar'}" {
		t.Errorf("data not as expected: %s", json)
	}

}
