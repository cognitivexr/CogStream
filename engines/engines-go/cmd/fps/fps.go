package main

import (
	"cognitivexr.at/cogstream/engines/pkg/engines/display"
	"cognitivexr.at/cogstream/pkg/webrtc"
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	err := webrtc.ServeEngineNetwork(ctx, display.Factory)
	log.Printf("engine server returned: %v", err)
}
