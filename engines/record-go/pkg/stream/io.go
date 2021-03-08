package stream

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const DefaultBufferSize = 1.5e+6 // 1.5MB

type FrameScanner struct {
	ctx StreamContext
	r   io.Reader
	lr  *io.LimitedReader

	buf   *bytes.Buffer
	frame []byte

	done bool
	err  error
}

func NewFrameScanner(ctx StreamContext, r io.Reader) *FrameScanner {
	buf := bytes.NewBuffer(make([]byte, DefaultBufferSize))
	buf.Reset()

	return &FrameScanner{
		ctx:  ctx,
		r:    r,
		lr:   &io.LimitedReader{R: r, N: 0},
		buf:  buf,
		done: false,
	}
}

func (s *FrameScanner) Next() bool {
	packet, err := s.readPacket()

	if err != nil {
		s.err = err
		s.done = true
		return false
	}

	s.frame = packet
	return true
}

func (s *FrameScanner) Err() error {
	return s.err
}

func (s *FrameScanner) readPacket() (data []byte, err error) {
	bufInt := make([]byte, 4)
	if _, err = s.r.Read(bufInt); err != nil {
		return
	}
	n := int64(binary.LittleEndian.Uint32(bufInt))

	// prepare limited reader to read the packet length exactly from the underlying reader
	if s.buf.Len() != 0 {
		panic(fmt.Sprintf("packet buffer should be empty, was %d", s.buf.Len()))
	}
	s.lr.N = n

	// prepare buffer to read into
	s.buf.Reset()
	s.buf.Grow(int(n)) // make sure we have enough space

	// read packet data into buffer
	_, err = s.buf.ReadFrom(s.lr)
	if err != nil {
		return
	}

	data = make([]byte, n)
	_, err = s.buf.Read(data)
	return
}
