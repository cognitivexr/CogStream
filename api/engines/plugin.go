package engines

import "context"

type PluginEngineRunner interface {
	Run(ctx context.Context, address string, specification *Specification) error
}
