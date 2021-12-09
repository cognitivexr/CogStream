package engines

import (
	"cognitivexr.at/cogstream/api/messages"
	"context"
)

type PluginEngineRunner interface {
	Run(ctx context.Context, startupObserver chan<- messages.EngineAddress, op messages.OperationSpec) error
}
