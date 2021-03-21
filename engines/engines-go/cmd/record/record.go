package main

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/format"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/engines/pkg/record"
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
	record.Serve(ctx, engine, "tcp", "0.0.0.0:53210")
}
