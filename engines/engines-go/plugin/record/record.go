package main

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/engines/pkg/engines/recorder"
	"cognitivexr.at/cogstream/pkg/serve"
	"context"
)

type runner struct {
}

func (r *runner) Run(ctx context.Context, op messages.OperationSpec, address string) error {
	err := serve.ServeEngineNetworkSingle(ctx, "tcp", address, recorder.Factory)
	return err
}

var Runner engines.PluginEngineRunner = &runner{}
