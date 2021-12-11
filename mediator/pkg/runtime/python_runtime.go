package runtime

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/mediator/pkg/runtime/cogstreampy"
)

func NewPythonRunner(modulePath string, descriptor *engines.EngineDescriptor) engines.PluginEngineRunner {
	return cogstreampy.NewEngineRunner(modulePath, descriptor)
}
