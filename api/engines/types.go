package engines

import (
	"cognitivexr.at/cogstream/api/messages"
	"time"
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

type InputFormat struct {
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	ColorMode ColorMode `json:"colorMode"`
}

// AnyInputFormat indicates that the any input format is supported or accepted.
var AnyInputFormat InputFormat

// Specification describes what an engine can do
type Specification struct {
	Operation   messages.OperationCode `json:"operation"`
	InputFormat InputFormat            `json:"inputFormat"`
	Attributes  messages.Attributes    `json:"attributes"`
}

type Engine struct {
	Name          string        `json:"name"`
	Specification Specification `json:"specification"`
}

type RunningEngine struct {
	Engine
	RuntimeId string
	Address   messages.EngineAddress
	Started   time.Time
	Stopped   time.Time
}
