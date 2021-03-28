package engines

import (
	"cognitivexr.at/cogstream/api/messages"
	"context"
)

type PluginEngineRunner interface {
	Run(ctx context.Context, op messages.OperationSpec, address string) error
}
