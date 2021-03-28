package engine

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/pkg/pipeline"
)

type Factory interface {
	Descriptor() engines.EngineDescriptor
	NewEngine() pipeline.Engine
}
