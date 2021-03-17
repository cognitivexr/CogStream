package engine

import (
	"cognitivexr.at/cogstream/api/format"
	"cognitivexr.at/cogstream/api/messages"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type contextKey struct {
	name string
}

var (
	StreamMetadataKey = &contextKey{name: "stream-metadata"}
)

type StreamMetadata struct {
	spec *messages.StreamSpec

	ClientFormat format.Format
	EngineFormat format.Format
}

func NewStreamMetadata() *StreamMetadata {
	return new(StreamMetadata)
}

func GetStreamMetadata(ctx context.Context) (metadata *StreamMetadata, ok bool) {
	metadata, ok = ctx.Value(StreamMetadataKey).(*StreamMetadata)
	return
}

func GetAttributes(ctx context.Context) (attr messages.Attributes, ok bool) {
	attr, ok = ctx.Value("attributes").(messages.Attributes)
	return
}

func NewStreamContext(parent context.Context, metadata *StreamMetadata) context.Context {
	return context.WithValue(parent, StreamMetadataKey, metadata)
}

func InitStream(ctx context.Context, r io.Reader) error {
	// read packet header
	n, err := readPacketHeader(r)
	if err != nil {
		return err
	}

	// de-serialize stream spec json
	lr := io.LimitReader(r, n)
	spec := new(messages.StreamSpec)
	err = json.NewDecoder(lr).Decode(spec)
	if err != nil {
		return fmt.Errorf("unable to decode stream spec: %v", err)
	}
	log.Printf("deserialized StreamSpec: %v", spec)

	clientFormat, err := messages.FormatFromAttributes(spec.Attributes)
	if err != nil {
		return fmt.Errorf("unable to determine client input format: %v", err)
	}

	metadata, ok := GetStreamMetadata(ctx)
	if ok && metadata != nil {
		metadata.ClientFormat = clientFormat
		log.Printf("setting ClientFormat in metadata %v", clientFormat)
	} else {
		log.Printf("could not set ClientFormat in metadata %v", clientFormat)
	}

	return nil
}
