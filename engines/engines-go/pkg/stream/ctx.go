package stream

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
	MetadataKey = &contextKey{name: "stream-metadata"}
)

type Metadata struct {
	StreamSpec *messages.StreamSpec

	ClientFormat format.Format
	EngineFormat format.Format
}

func NewMetadata() *Metadata {
	return new(Metadata)
}

func GetMetadata(ctx context.Context) (metadata *Metadata, ok bool) {
	metadata, ok = ctx.Value(MetadataKey).(*Metadata)
	return
}

func GetAttributes(ctx context.Context) (attr messages.Attributes, ok bool) {
	attr, ok = ctx.Value("attributes").(messages.Attributes)
	return
}

func ContextWithMetadata(parent context.Context, metadata *Metadata) context.Context {
	return context.WithValue(parent, MetadataKey, metadata)
}

func InitStream(ctx context.Context, r io.Reader) error {
	// read packet header
	n, err := readInt(r)
	if err != nil {
		return err
	}

	// de-serialize stream StreamSpec json
	lr := io.LimitReader(r, n)
	spec := new(messages.StreamSpec)
	err = json.NewDecoder(lr).Decode(spec)
	if err != nil {
		return fmt.Errorf("unable to decode stream StreamSpec: %v", err)
	}
	log.Printf("deserialized StreamSpec: %v", spec)

	clientFormat, err := messages.FormatFromAttributes(spec.Attributes)
	if err != nil {
		return fmt.Errorf("unable to determine client input format: %v", err)
	}

	metadata, ok := GetMetadata(ctx)
	if ok && metadata != nil {
		metadata.ClientFormat = clientFormat
		log.Printf("setting ClientFormat in metadata %v", clientFormat)
	} else {
		log.Printf("could not set ClientFormat in metadata %v", clientFormat)
	}

	return nil
}
