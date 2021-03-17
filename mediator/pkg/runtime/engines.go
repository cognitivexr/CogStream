package runtime

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
)

type EngineFinder interface {
	ListEngines() []*engines.Engine
	FindEngines(engines.Specification) []*engines.Engine
	FindEngineByName(string) (*engines.Engine, bool)
}

type EngineRuntime interface {
	StartEngine(engine *engines.Engine, attributes messages.Attributes) (*engines.RunningEngine, error)
	StopEngine(*engines.RunningEngine) error
	ListRunning() []*engines.RunningEngine
}
