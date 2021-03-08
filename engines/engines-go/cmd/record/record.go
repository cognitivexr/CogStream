package main

import (
	"cognitivexr.at/cogstream/engines/pkg/record"
)

func main() {
	record.Serve("tcp", "0.0.0.0:53210")
}
