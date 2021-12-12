package rpc

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"testing"
	"time"
)

type DummyRuntime struct {
}

func (d DummyRuntime) StartEngine(engine *engines.EngineDescriptor, spec messages.OperationSpec) (*engines.RunningEngine, error) {
	log.Info(">> START ENGINE")
	return nil, nil
}

func (d DummyRuntime) StopEngine(engine *engines.RunningEngine) error {
	log.Info(">> STOP ENGINE")
	return nil
}

func (d DummyRuntime) ListRunning() []*engines.RunningEngine {
	//TODO implement me
	panic("implement me")
}

func TestName(t *testing.T) {
	remote := NewEngineRuntimeSkeleton(&DummyRuntime{})
	go remote.Serve(":45341")

	time.Sleep(1 * time.Second)
	log.Info("Invoking skeleton")
	stub := NewEngineRuntimeStub("127.0.0.1:45341")
	_, err := stub.StartEngine(&engines.EngineDescriptor{}, messages.OperationSpec{})
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	time.Sleep(1 * time.Second)

	log.Info("ok bye!")
	remote.Close()
}
