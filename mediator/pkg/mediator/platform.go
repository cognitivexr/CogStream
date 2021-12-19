package mediator

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"cognitivexr.at/cogstream/mediator/pkg/runtime"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var EngineNotFound = errors.New("engine not found")
var EngineNotStarted = errors.New("engine not started")

// TODO: Alert when an error happens here
type Platform interface {
	GetAvailableEngines(*HandshakeContext) (*messages.AvailableEngines, error)
	GetStreamSpec(*HandshakeContext) (*messages.StreamSpec, error)
	ListAvailableEngines() (*messages.AvailableEngines, error)
}

type DummyPlatform struct {
	runtime runtime.EngineRuntime
	finder  runtime.EngineFinder
}

func NewPluginPlatform(pluginDirs ...string) (Platform, error) {
	engineRuntime := runtime.NewPluginEngineRuntime(pluginDirs...)
	err := engineRuntime.LoadPlugins()

	if err != nil {
		return nil, err
	}

	return &DummyPlatform{
		runtime: engineRuntime,
		finder:  engineRuntime,
	}, nil
}

func (d *DummyPlatform) ListAvailableEngines() (*messages.AvailableEngines, error) {
	foundEngines := d.finder.ListEngines()

	availableEngines := &messages.AvailableEngines{Engines: make([]*messages.EngineSpec, 0)}

	for _, engine := range foundEngines {
		spec := &messages.EngineSpec{Attributes: messages.NewAttributes()}
		mapAvailableEngines(engine, spec)
		availableEngines.Engines = append(availableEngines.Engines, spec)
	}

	return availableEngines, nil
}

func (d *DummyPlatform) GetAvailableEngines(hs *HandshakeContext) (*messages.AvailableEngines, error) {
	foundEngines := d.finder.FindEngines(engines.Specification{
		Operation:  hs.OperationSpec.Code,
		Attributes: hs.OperationSpec.Attributes,
	})

	availableEngines := &messages.AvailableEngines{Engines: make([]*messages.EngineSpec, 0)}

	for _, engine := range foundEngines {
		spec := &messages.EngineSpec{Attributes: messages.NewAttributes()}
		mapAvailableEngines(engine, spec)
		availableEngines.Engines = append(availableEngines.Engines, spec)
	}

	return availableEngines, nil
}

func mapAvailableEngines(engine *engines.EngineDescriptor, engineSpec *messages.EngineSpec) {
	engineSpec.Name = engine.Name
	engineSpec.Attributes.Set("format.width", strconv.Itoa(engine.Specification.InputFormat.Width))
	engineSpec.Attributes.Set("format.height", strconv.Itoa(engine.Specification.InputFormat.Height))
	engineSpec.Attributes.Set("format.colorMode", strconv.Itoa(int(engine.Specification.InputFormat.ColorMode)))
	for k, v := range engine.Specification.Attributes {
		engineSpec.Attributes[k] = v
	}
}

func (d *DummyPlatform) GetStreamSpec(hs *HandshakeContext) (*messages.StreamSpec, error) {
	engine, ok := d.finder.FindEngineByName(hs.ClientFormatSpec.Engine)
	if !ok {
		// TODO: alert
		return nil, errors.New("cannot find engine")
	}
	log.Info("found engine %v", engine)

	runningEngine, err := d.runtime.StartEngine(engine, *hs.OperationSpec)
	if err != nil {
		// TODO: alert that runtime f'd up
		log.Warn("engine %s failed to start: %v", engine.Name, err)
		return nil, EngineNotStarted
	}
	if runningEngine == nil {
		return nil, fmt.Errorf("engine runtime returned nil after StartEngine")
	}

	// TODO stream format should be negotiated, here we just assume that the stream format = client format
	attrs := messages.NewAttributes()
	format, err := messages.FormatFromAttributes(hs.ClientFormatSpec.Attributes)
	if err == nil {
		messages.FormatToAttributes(format, attrs)
	} else {
		log.Error("failed to read client format from attributes", err)
	}

	// TODO:
	address := strings.Replace(string(runningEngine.Address), "0.0.0.0", "127.0.0.1", 1)
	return &messages.StreamSpec{EngineAddress: messages.EngineAddress(address), Attributes: attrs}, nil
}
