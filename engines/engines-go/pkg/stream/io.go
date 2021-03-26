package stream

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
)

const DefaultBufferSize = 1.5e+6 // 1.5MB
const FramePacketHeaderSize = 24

type FramePacketHeader struct {
	StreamId         uint32
	FrameId          uint32
	TimestampSeconds uint32
	TimestampNanos   uint32
	MetadataLen      uint32
	DataLen          uint32
}

type FramePacket struct {
	Header   FramePacketHeader
	Metadata []byte
	Data     []byte
}

// FramePacketReader reads FramePacket instances from an underlying stream.
type FramePacketReader interface {
	ReadPacket(packet *FramePacket) error
}

func NewFramePacketReader(r io.Reader) FramePacketReader {
	bBuf := bytes.NewBuffer(make([]byte, DefaultBufferSize))
	bBuf.Reset()

	hBuf := bytes.NewBuffer(make([]byte, FramePacketHeaderSize))
	hBuf.Reset()

	return &framePacketReader{
		r:    r,
		bLr:  &io.LimitedReader{R: r, N: 0},
		hLr:  &io.LimitedReader{R: r, N: FramePacketHeaderSize},
		hBuf: hBuf,
		bBuf: bBuf,
	}
}

type framePacketReader struct {
	r io.Reader // r holds the underlying reader

	hLr  *io.LimitedReader // hLr is the limited reader for the packet header
	hBuf *bytes.Buffer     // hBuf is a buffer to read the packet header
	bLr  *io.LimitedReader // bLr is the limited reader for the packet body
	bBuf *bytes.Buffer     // bBuf is a buffer to read the packet body
}

func (s *framePacketReader) ReadPacket(packet *FramePacket) error {
	err := s.readPacketHeader(&packet.Header)
	if err != nil {
		return err
	}

	n := int(packet.Header.DataLen + packet.Header.MetadataLen)

	// prepare limited reader to read the packet length exactly from the underlying reader
	buf := s.bBuf
	if buf.Len() != 0 {
		panic(fmt.Sprintf("packet buffer should be empty, was %d", buf.Len()))
	}
	s.bLr.N = int64(n)

	// prepare buffer to read into
	buf.Reset()
	buf.Grow(n) // make sure we have enough space

	// read packet data into buffer
	read, err := buf.ReadFrom(s.bLr)
	if err != nil {
		return err
	}
	if read != int64(n) {
		return fmt.Errorf("expected %d packet bytes, read %d", n, read)
	}

	// read optional metadata
	if packet.Header.MetadataLen > 0 {
		packet.Metadata = make([]byte, packet.Header.MetadataLen)
		n, err = buf.Read(packet.Metadata)
		if err != nil {
			return err
		}
		if n != int(packet.Header.MetadataLen) {
			return fmt.Errorf("expected %d metadata bytes, received %d", packet.Header.MetadataLen, n)
		}
	}

	// read data
	packet.Data = make([]byte, packet.Header.DataLen)
	n, err = buf.Read(packet.Data)
	if err != nil {
		return err
	}
	if n != int(packet.Header.DataLen) {
		return fmt.Errorf("expected %d data bytes, received %d", packet.Header.DataLen, n)
	}

	return nil
}

func (s *framePacketReader) readPacketHeader(header *FramePacketHeader) error {
	buf := s.hBuf
	buf.Reset()
	s.hLr.N = FramePacketHeaderSize

	n, err := buf.ReadFrom(s.hLr)
	if err != nil {
		return err
	}
	if n != FramePacketHeaderSize {
		if n == 0 {
			return io.EOF
		}
		return fmt.Errorf("expected %d header bytes, read %d", FramePacketHeaderSize, n)
	}

	header.StreamId = nextUint32(buf)
	header.FrameId = nextUint32(buf)
	header.TimestampSeconds = nextUint32(buf)
	header.TimestampNanos = nextUint32(buf)
	header.MetadataLen = nextUint32(buf)
	header.DataLen = nextUint32(buf)

	return nil
}

// FramePacketScanner uses a FramePacketReader to create iterator-like interface and takes care of memory allocation.
type FramePacketScanner struct {
	r FramePacketReader

	frame *FramePacket // last read frame

	done bool
	err  error
}

func NewFramePacketScanner(r FramePacketReader) *FramePacketScanner {
	return &FramePacketScanner{
		r:    r,
		done: false,
	}
}

func (s *FramePacketScanner) Scan(ctx context.Context) (*FramePacket, error) {
	if s.Next() {
		return s.frame, nil
	} else {
		return nil, s.Err()
	}
}

func (s *FramePacketScanner) Next() bool {
	var packet FramePacket
	err := s.r.ReadPacket(&packet)

	if err != nil {
		s.err = err
		s.done = true
		return false
	}

	s.frame = &packet
	return true
}

func (s *FramePacketScanner) Get() *FramePacket {
	return s.frame
}

func (s *FramePacketScanner) Err() error {
	return s.err
}

func readInt(r io.Reader) (int64, error) {
	bufInt := make([]byte, 4)
	if _, err := r.Read(bufInt); err != nil {
		return -1, err
	}
	n := binary.LittleEndian.Uint32(bufInt)
	return int64(n), nil
}

func nextUint32(buf *bytes.Buffer) uint32 {
	return binary.LittleEndian.Uint32(buf.Next(4))
}
