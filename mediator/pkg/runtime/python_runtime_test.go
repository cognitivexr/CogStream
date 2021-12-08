package runtime

import (
	"cognitivexr.at/cogstream/api/messages"
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

	go func() {
		runner.Run(ctx, messages.OperationSpec{}, "0.0.0.0:45312")
		wg.Done()
	}()

	time.Sleep(10 * time.Second)
	cancel()
	wg.Wait()
}
