package main

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/engines/pkg/engines/recorder"
	"context"
)

type runner struct {
}

func (r *runner) Run(ctx context.Context, engine *engines.EngineDescriptor, address string) error {
	err := recorder.ServeSingle(ctx, engine, "tcp", address)
	return err
}

var Runner engines.PluginEngineRunner = &runner{}
