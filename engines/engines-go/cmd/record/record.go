package main

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/format"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/engines/pkg/engines/recorder"
	"context"
)

func main() {
	ctx := context.Background()

	engine := &engines.EngineDescriptor{
		Name: "record",
		Specification: engines.Specification{
			Operation:   messages.OperationRecord,
			InputFormat: format.AnyFormat,
			Attributes:  messages.NewAttributes(),
		},
	}
	recorder.Serve(ctx, engine, "tcp", "0.0.0.0:53210")
}
