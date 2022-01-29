package runtime

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/mediator/pkg/runtime/docker"
)

func NewDockerRunner(modulePath string, descriptor *engines.EngineDescriptor) engines.PluginEngineRunner {
	return docker.NewEngineRunner(modulePath, descriptor)
}
