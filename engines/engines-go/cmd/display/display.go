package main

import (
	"cognitivexr.at/cogstream/engines/pkg/engines/display"
	"cognitivexr.at/cogstream/pkg/serve"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	err := serve.ServeEngineNetwork(ctx, "tcp", "0.0.0.0:54321", display.Factory)
	log.Printf("engine server returned: %v", err)
}
