package main

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/engines/pkg/record"
)

type runner struct {
}

func (r *runner) Run(address string, specification *engines.Specification) error {
	err := record.ServeSingle("tcp", address)
	return err
}

var Runner engines.PluginEngineRunner = &runner{}
