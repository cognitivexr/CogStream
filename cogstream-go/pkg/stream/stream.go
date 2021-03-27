package stream

import (
	"cognitivexr.at/cogstream/api/messages"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Stream struct {
	io.ReadWriter

	metadata *Metadata
	ctx      context.Context
}

type Configurator func(stream *Stream)

func (s *Stream) Metadata() *Metadata {
	return s.metadata
}

func NewStream(ctx context.Context, rw io.ReadWriter) (*Stream, error) {
	// read stream spec
	spec := new(messages.StreamSpec)
	err := ReadStreamSpec(rw, spec)
	if err != nil {
		return nil, fmt.Errorf("error reading stream StreamSpec: %v", err)
	}

	// initialize metadata
	metadata := NewMetadata()
	metadata.StreamSpec = spec

	// read client format from stream spec
	if format, err := messages.FormatFromAttributes(spec.Attributes); err == nil {
		metadata.ClientFormat = format
	} else {
		return nil, err
	}

	// initialize stream context (context with Metadata)
	ctx = ContextWithMetadata(ctx, metadata)

	return &Stream{rw, metadata, ctx}, nil
}

func ReadStreamSpec(r io.Reader, spec *messages.StreamSpec) error {
	// read packet header
	n, err := readInt(r)
	if err != nil {
		return fmt.Errorf("error reading StreamSpec header: %v", err)
	}

	// de-serialize stream StreamSpec json
	lr := io.LimitReader(r, n)
	err = json.NewDecoder(lr).Decode(spec)

	if err != nil {
		return fmt.Errorf("error decoding JSON StreamSpec: %v", err)
	}

	log.Printf("deserialized StreamSpec: %v", spec)
	return nil
}
