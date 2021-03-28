package main

import (
	"cognitivexr.at/cogstream/engines/pkg/engines/recorder"
	"cognitivexr.at/cogstream/pkg/serve"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	err := serve.ServeEngineNetwork(ctx, "tcp", "0.0.0.0:53210", recorder.Factory)
	log.Printf("engine server returned: %v", err)
}
