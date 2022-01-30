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

type EngineDescriptor struct {
	Name          string            `json:"name"`
	Runtime       string            `json:"runtime"`
	RuntimeConfig map[string]string `json:"runtime-config"`
	Specification Specification     `json:"specification"`
}

type RunningEngine struct {
	EngineDescriptor
	RuntimeId string                 `json:"runtimeId"`
	Address   messages.EngineAddress `json:"address"`
	Started   time.Time              `json:"started"`
	Stopped   time.Time              `json:"stopped"`
}
