package engine

import (
	"context"
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

type Config struct {
	Width, Height int
	ColorMode     ColorMode
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
