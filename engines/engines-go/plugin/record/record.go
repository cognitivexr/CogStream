package main

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/engines/pkg/record"
	"context"
)

type runner struct {
}

func (r *runner) Run(ctx context.Context, engine *engines.Engine, address string) error {
	err := record.ServeSingle(ctx, engine, "tcp", address)
	return err
}

var Runner engines.PluginEngineRunner = &runner{}
