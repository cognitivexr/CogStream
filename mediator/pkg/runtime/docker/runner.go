package docker

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"context"
)

type EngineRunner struct {
	modulePath string
	descriptor *engines.EngineDescriptor
}

// NewEngineRunner creates a new EngineRunner that creates python engine processes for each Run invocation.
func NewEngineRunner(modulePath string, descriptor *engines.EngineDescriptor) *EngineRunner {
	return &EngineRunner{
		modulePath: modulePath,
		descriptor: descriptor,
	}
}

func (p *EngineRunner) Run(ctx context.Context, startupObserver chan<- messages.EngineAddress, op messages.OperationSpec) error {
	return NewEngineContainer(p.modulePath, p.descriptor).Run(ctx, startupObserver, op)
}
