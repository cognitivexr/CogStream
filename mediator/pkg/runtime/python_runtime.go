package runtime

import (
	"bufio"
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"context"
	"errors"
	"fmt"
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
	mutex      sync.Mutex
}

func NewPythonRunner(modulePath string, descriptor *engines.EngineDescriptor) *PythonRunner {
	return &PythonRunner{
		modulePath: modulePath,
		descriptor: descriptor,
	}
}

func parseEngineAddressFromLog(line string) (string, error) {
	// INFO:cogstream.engine.srv:started server socket on address ('0.0.0.0', 46699)
	tupleStr := strings.Trim(line[59:], " \n\r()")
	tuple := strings.Split(tupleStr, ",")
	if len(tuple) != 2 {
		return "", errors.New("tuple string did not contain two elements: " + tupleStr)
	}
	addrStr, portStr := tuple[0], tuple[1]
	addrStr = strings.Trim(addrStr, "'")
	portStr = strings.TrimSpace(portStr)

	return fmt.Sprintf("%s:%s", addrStr, portStr), nil
}

func (p *PythonRunner) Run(ctx context.Context, startupObserver chan<- messages.EngineAddress, op messages.OperationSpec) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.cmd != nil {
		return nil
	}

	// if the plugin is self-contained in a venv then we use the python3 bin from the venv
	pythonCommand := path.Join(p.modulePath, ".venv", "bin", "python3")

	if _, err := os.Stat(pythonCommand); errors.Is(err, os.ErrNotExist) {
		// otherwise we'll use the system-wide python binary
		pythonCommand = "/usr/bin/python3"
	}

	// we assume the plugin has a `main` module inside the plugin module
	module := p.descriptor.Name + ".main"

	cmd := exec.Command(pythonCommand, "-m", module, "--host", "0.0.0.0", "--port", "0")
	p.cmd = cmd
	cmd.Dir = p.modulePath // chdir
	cmd.Stdout = os.Stdout
	console, err := cmd.StderrPipe()

	notifyStarted := func() {
		log.Info("started engine serving")
		if started, ok := ctx.Value("started").(*sync.WaitGroup); started != nil && ok {
			started.Done()
		}
	}

	notifyAddress := func(str string) {
		startupObserver <- messages.EngineAddress(str)
	}

	startupDone := false

	// startup/log watchdog
	go func() {
		reader := bufio.NewReader(console)
		line, err := reader.ReadString('\n')
		for err == nil {
			line = strings.TrimRight(line, " \n\r")
			fmt.Println(line) // forward to stdout

			if !startupDone {
				if strings.HasPrefix(line, "INFO:cogstream.engine.srv:started server socket on address") {
					addr, parseErr := parseEngineAddressFromLog(line)
					if parseErr != nil {
						log.Error("could not determine engine address for from line `%s`", line)
						continue // TODO: engine should fail
					}
					notifyAddress(addr)
				}
				// mark engine as started when we receive the respective log output from the engine
				if line == "INFO:cogstream.engine.srv:waiting for next connection" {
					notifyStarted()
					startupDone = true
				}
			}

			line, err = reader.ReadString('\n')
		}
	}()

	// start kill watchdog
	go func() {
		select {
		case <-ctx.Done():
			cmd.Process.Kill()
		}
	}()

	// run command
	log.Info("starting command %s", cmd)
	err = cmd.Start()
	if err != nil {
		return err
	}

	// block
	return cmd.Wait()
}
