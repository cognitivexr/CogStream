package main

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/engines/pkg/record"
	"context"
)

type runner struct {
}

func (r *runner) Run(ctx context.Context, address string, specification *engines.Specification) error {
	err := record.ServeSingle(ctx,"tcp", address)
	return err
}

var Runner engines.PluginEngineRunner = &runner{}
