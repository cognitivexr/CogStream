package engines

import "context"

type PluginEngineRunner interface {
	Run(ctx context.Context, engine *EngineDescriptor, address string) error
}
