package main

import (
	"cognitivexr.at/cogstream/engines/pkg/engines/stats"
	"cognitivexr.at/cogstream/pkg/webrtc"
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	err := webrtc.ServeEngineNetwork(ctx, stats.Factory)
	log.Printf("engine server returned: %v", err)
}
