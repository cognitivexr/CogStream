package runtime

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"context"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

type PythonRunner struct {
	modulePath string
	descriptor *engines.EngineDescriptor
	cmd        *exec.Cmd
}

func NewPythonRunner(modulePath string, descriptor *engines.EngineDescriptor) *PythonRunner {
	return &PythonRunner{
		modulePath: modulePath,
		descriptor: descriptor,
	}
}

func (p *PythonRunner) Run(ctx context.Context, op messages.OperationSpec, address string) error {
	// TODO: lock
	if p.cmd != nil {
		return nil
	}

	parts := strings.Split(address, ":")
	host, port := parts[0], parts[1]

	// we assume the plugin is self-contained in a venv
	pythonCommand := path.Join(p.modulePath, ".venv", "bin", "python3")
	// and that the plugin has a main module
	module := p.descriptor.Name + ".main"

	cmd := exec.Command(pythonCommand, "-m", module, "--host", host, "--port", port)
	p.cmd = cmd
	cmd.Dir = p.modulePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Info("starting command", cmd.Start())

	if started, ok := ctx.Value("started").(*sync.WaitGroup); started != nil && ok {
		log.Println("started engine serving")
		started.Done()
	}

	go func() {
		select {
		case <-ctx.Done():
			cmd.Process.Kill()
		}
	}()

	return cmd.Wait()
}
