package engines

import (
	"cognitivexr.at/cogstream/api/format"
	"cognitivexr.at/cogstream/api/messages"
	"time"
)

// Specification describes what an engine can do
type Specification struct {
	Operation   messages.OperationCode `json:"operation"`
	InputFormat format.Format          `json:"inputFormat"`
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
