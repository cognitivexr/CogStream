package mediator

import (
	"cognitivexr.at/cogstream/api/messages"
	"log"
	"testing"
)

// TODO: extend tests

func TestMarshalEngineFormatSpec(t *testing.T) {
	efs := messages.AvailableEngines{}
	marshalled, err := MarshalAvailableEngines(&efs)
	if err != nil {
		t.Errorf("failed to marshal: %v", err)
	}
	log.Printf("result: %s", marshalled)

}
