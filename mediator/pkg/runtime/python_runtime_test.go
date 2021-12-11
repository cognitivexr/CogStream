package runtime

import (
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"context"
	"path"
	"sync"
	"testing"
	"time"
)

func TestRunPlugin(t *testing.T) {
	t.Skip("this test exists for development purposes")

	dir := "/home/thomas/workspace/cognitivexr/cogstream/engines/engines-py/fermx"
	spec, _ := ParseSpec(path.Join(dir, "fermx.spec.json"))
	runner := NewPythonRunner(dir, spec)
	ctx, cancel := context.WithCancel(context.TODO())

	wg := sync.WaitGroup{}
	wg.Add(1)
	obs := make(chan messages.EngineAddress, 1)

	go func() {
		runner.Run(ctx, obs, messages.OperationSpec{})
		wg.Done()
	}()

	addr := <-obs
	log.Info("engine address is %s", addr)

	time.Sleep(10 * time.Second)
	cancel()
	wg.Wait()
}
