package mediator

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/pkg/log"
	"cognitivexr.at/cogstream/pkg/runtime"
	"fmt"
	"time"
)

type Platform interface {
	GetEngineFormatSpec(*HandshakeContext) (*messages.EngineFormatSpec, error)
	GetStreamSpec(*HandshakeContext) (*messages.StreamSpec, error)
}

type DummyPlatform struct {
	runtime runtime.EngineRuntime
	finder  runtime.EngineFinder
}

func NewPluginPlatform(pluginDir string) (Platform, error) {
	engineRuntime := runtime.NewPluginEngineRuntime(pluginDir)
	err := engineRuntime.LoadPlugins()

	if err != nil {
		return nil, err
	}

	return &DummyPlatform{
		runtime: engineRuntime,
		finder:  engineRuntime,
	}, nil
}

func (d *DummyPlatform) GetEngineFormatSpec(_ *HandshakeContext) (*messages.EngineFormatSpec, error) {
	attributes := make(messages.Attributes)
	attributes.Set("foo", "bar")
	return &messages.EngineFormatSpec{Attributes: attributes}, nil
}

func (d *DummyPlatform) GetStreamSpec(hs *HandshakeContext) (*messages.StreamSpec, error) {
	attributes := make(messages.Attributes)
	attributes.Set("sessionId", hs.SessionId)
	attributes.Set("config.width", "640")
	attributes.Set("config.height", "360")
	attributes.Set("config.colormode", "1")

	var address messages.EngineAddress = "127.0.0.1:53210"

	availableEngines := d.finder.FindEngines(engines.Specification{Operation: hs.OperationSpec.Code})

	if len(availableEngines) <= 0 {
		return nil, fmt.Errorf("no engine for the given operation")
	}
	if len(availableEngines) > 1 {
		// TODO: decide how to resolve ambiguities
		log.Warn("too many engines match the specification")
	}

	engine := availableEngines[0]
	log.Info("found engine %v", engine)

	go func() {
		_, err := d.runtime.StartEngine(engine)
		if err != nil {
			log.Warn("engine %s stopped with error: %s", engine.Name, err)
		} else {
			log.Info("engine %s stopped", engine.Name)
		}
	}()

	time.Sleep(500 * time.Millisecond) // TODO: how to indicate that the engine is up and running?

	return &messages.StreamSpec{EngineAddress: address, Attributes: attributes}, nil
}
