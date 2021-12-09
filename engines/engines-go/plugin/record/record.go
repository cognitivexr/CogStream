package main

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/engines/pkg/engines/recorder"
	"cognitivexr.at/cogstream/pkg/serve"
	"context"
	"fmt"
	"net"
)

type runner struct {
}

func getRandomTpcPort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	return listener.Addr().(*net.TCPAddr).Port
}

func (r *runner) Run(ctx context.Context, startupObserver chan<- messages.EngineAddress, op messages.OperationSpec) error {
	addr := fmt.Sprintf("0:0:0:0:%d", getRandomTpcPort())
	startupObserver <- messages.EngineAddress(addr)
	err := serve.ServeEngineNetworkSingle(ctx, "tcp", addr, recorder.Factory)
	return err
}

var Runner engines.PluginEngineRunner = &runner{}
