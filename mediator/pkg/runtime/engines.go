package runtime

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
)

type EngineFinder interface {
	ListEngines() []*engines.EngineDescriptor
	FindEngines(engines.Specification) []*engines.EngineDescriptor
	FindEngineByName(string) (*engines.EngineDescriptor, bool)
}

type EngineRuntime interface {
	StartEngine(engine *engines.EngineDescriptor, spec messages.OperationSpec) (*engines.RunningEngine, error)
	StopEngine(*engines.RunningEngine) error
	ListRunning() []*engines.RunningEngine
}
