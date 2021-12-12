package main

import (
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"flag"
)

func main() {
	hostPtr := flag.String("host", "0.0.0.0", "host to bind to")
	portPtr := flag.Int("port", 8191, "the server port")
	flag.Parse()

	log.Info("starting cluster leader on %s:%d", hostPtr, portPtr)
}
