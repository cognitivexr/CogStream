package docker

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
	"strings"
	"sync"
)

//TODO use with docker log
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

type engineContainer struct {
	image      string
	descriptor *engines.EngineDescriptor
	container  int // TODO use docker API
	mutex      sync.Mutex
	exit       sync.WaitGroup
	address    string // holds the server address and port if the plugin is running
	err        error
}

func NewEngineContainer(modulePath string, descriptor *engines.EngineDescriptor) *engineContainer {
	image := descriptor.RuntimeConfig.Get("image")
	ct := &engineContainer{
		image:      image,
		descriptor: descriptor,
	}
	ct.exit.Add(1)
	return ct
}

func (c *engineContainer) Interrupt() error {
	return nil
}

func (c *engineContainer) Wait() error {
	return nil
}

func (c *engineContainer) Run(ctx context.Context, startupObserver chan<- messages.EngineAddress, op messages.OperationSpec) error {
	c.mutex.Lock()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, c.image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	_, err = cli.ContainerCreate(ctx, &container.Config{
		Image: c.image,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	// block until process is done
	c.mutex.Unlock()
	// notify exit
	c.exit.Done()
	return nil
}

func (c *engineContainer) notifyStarted(ctx context.Context) {
	log.Info("started engine serving")
	if started, ok := ctx.Value("started").(*sync.WaitGroup); started != nil && ok {
		started.Done()
	}
}

func (c *engineContainer) notifyAddress(startupObserver chan<- messages.EngineAddress) {
	startupObserver <- messages.EngineAddress(c.address)
}
