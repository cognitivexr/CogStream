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
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func getFreeTcpPort() int {
	var a *net.TCPAddr
	var err error
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port
		}
	}

	panic(err)
}

type engineContainer struct {
	image      string
	descriptor *engines.EngineDescriptor
	container  string
	mutex      sync.Mutex
	exit       sync.WaitGroup
	address    string // holds the server address and port if the plugin is running
	err        error
}

func NewEngineContainer(modulePath string, descriptor *engines.EngineDescriptor) *engineContainer {
	image := descriptor.RuntimeConfig["image"]
	ct := &engineContainer{
		image:      image,
		descriptor: descriptor,
	}
	ct.exit.Add(1)
	return ct
}

func (c *engineContainer) Interrupt() error {
	if c.container == "" {
		return errors.New("container has not yet started")
	}
	log.Info("interrupting engine at %s", c.address)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	return c.shutdown(cli)
}

func (c *engineContainer) shutdown(cli *client.Client) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.container == "" {
		return errors.New("container has not yet started")
	}
	ctx := context.Background()
	timeout, _ := time.ParseDuration("10s")
	err := cli.ContainerStop(ctx, c.container, &timeout)
	if err != nil {
		return err
	}
	err = cli.ContainerRemove(ctx, c.container, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	c.container = ""
	return nil
}

func (c *engineContainer) Wait() error {
	c.mutex.Lock()

	if c.container == "" {
		c.mutex.Unlock()
		return errors.New("container has not yet started")
	}

	c.mutex.Unlock()
	c.exit.Wait()

	return c.err
}

func checkHostAvailable(address string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return fmt.Errorf("engine not available: %s", err)
	}
	if conn != nil {
		defer conn.Close()
		log.Debug("engine available: %s", address)
	}
	// TODO needs to be removed still.
	time.Sleep(time.Second)
	return nil
}

func (c *engineContainer) Run(ctx context.Context, startupObserver chan<- messages.EngineAddress, op messages.OperationSpec) error {
	// TODO: containers are never removed

	c.mutex.Lock()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// require image
	var image string
	images, err := cli.ImageList(ctx, types.ImageListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("reference", c.image)),
	})
	if len(images) == 0 {
		image = c.image
		log.Info("could not found image %s, trying to pull", c.image)
		reader, err := cli.ImagePull(ctx, c.image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		io.Copy(os.Stdout, reader)
	} else if len(images) > 1 {
		image = images[0].ID
		// TODO: more than one image with the given reference was found?
	} else {
		image = images[0].ID
	}

	engineHost := "0.0.0.0"
	enginePort := strconv.Itoa(getFreeTcpPort())

	portBindings := nat.PortMap{
		// the default engine port is 54321 (exposed by Dockerfiles)
		nat.Port("54321/tcp"): []nat.PortBinding{
			{HostIP: engineHost, HostPort: enginePort},
		},
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: image,
		},
		&container.HostConfig{PortBindings: portBindings},
		nil,
		nil,
		"",
	)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	c.container = resp.ID
	c.address = engineHost + ":" + enginePort
	c.notifyAddress(startupObserver)

	err = checkHostAvailable(c.address, time.Minute)
	if err != nil {
		return fmt.Errorf("engine took too long to start: %s", err)
	}

	c.notifyStarted(ctx)

	c.mutex.Unlock()
	// block until container is done
	defer c.shutdown(cli)
	defer c.exit.Done()

	_, errC := cli.ContainerWait(ctx, c.container, "not-running")
	if err := <-errC; err != nil {
		// this may also be raised if the context is cancelled
		return err
	}
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
