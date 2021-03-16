package engines

import "context"

type PluginEngineRunner interface {
	Run(ctx context.Context, engine *Engine, address string) error
}
