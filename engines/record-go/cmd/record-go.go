package main

import (
	"cognitivexr.at/cogstream/engines/record-go/pkg/stream"
)

func main() {
	stream.Serve("tcp", "0.0.0.0:53210")
}
