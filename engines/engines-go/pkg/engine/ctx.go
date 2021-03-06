package engine

import (
	"context"
	"encoding/json"
	"io"
	"log"
)

type ColorMode int

const (
	UNKNOWN = iota
	RGB
	RGBA
	GRAY
	BGR
	BGRA
	HLS
	Lab
	Luv
	Bayer
)

// FIXME: this should be the same struct as the mediator provides
type Agreement struct {
	SessionId string `json:"sessionId"`
	Operation string `json:"operation"`
	Config    Config `json:"config"`
}

type Config struct {
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	ColorMode ColorMode `json:"colorMode"`
}

type StreamContext interface {
	context.Context
	Config() *Config
}

type defaultStreamContext struct {
	context.Context
	config *Config
}

func (d *defaultStreamContext) Config() *Config {
	return d.config
}

func NewStreamContext() *defaultStreamContext {
	return &defaultStreamContext{context.TODO(), new(Config)}
}

func InitStreamContext(r io.Reader) (StreamContext, error) {
	// read packet header
	n, err := readPacketHeader(r)
	if err != nil {
		return nil, err
	}

	// de-serialize json
	lr := io.LimitReader(r, n)
	a := new(Agreement)
	err = json.NewDecoder(lr).Decode(a)
	if err != nil {
		return nil, err
	}
	log.Printf("deserialized agreement %v\n", a)

	ctx := NewStreamContext()
	ctx.config = &a.Config
	return ctx, nil
}
