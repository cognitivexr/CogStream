package runtime

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"cognitivexr.at/cogstream/mediator/pkg/util"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
	"sync"
)

type PluginEngine struct {
	Path   string
	Runner engines.PluginEngineRunner
	Engine *engines.EngineDescriptor
}

func OpenPluginEngine(pluginPath string) (*PluginEngine, error) {
	specPath := pluginPath + ".spec.json"

	spec, err := ParseSpec(specPath)
	if err != nil {
		return nil, err
	}

	pluginObj, err := plugin.Open(pluginPath)

	if err != nil {
		return nil, err
	}

	symbol, err := pluginObj.Lookup("Runner")

	if err != nil {
		return nil, err
	}
	if symbol == nil {
		return nil, fmt.Errorf("symbol Runner of plugin %s was nil", spec.Name)
	}

	runner, ok := symbol.(*engines.PluginEngineRunner)
	if !ok {
		return nil, fmt.Errorf("runner of plugin %s should be a PluginEngineRunner, but is a %v", spec.Name, reflect.TypeOf(symbol))
	}

	return &PluginEngine{
		pluginPath,
		*runner,
		spec,
	}, nil
}

func CreatePythonPluginEngine(descriptorPath string, descriptor *engines.EngineDescriptor) (*PluginEngine, error) {
	pluginPath := path.Dir(descriptorPath)
	log.Info("loading python plugin engine in %s", pluginPath)

	return &PluginEngine{
		pluginPath,
		NewPythonRunner(pluginPath, descriptor),
		descriptor,
	}, nil
}

type runningEngineContext struct {
	runningEngine *engines.RunningEngine
	ctx           context.Context
	cancel        context.CancelFunc
	started       *sync.WaitGroup
}

type pluginEngineRuntime struct {
	pluginDirs       []string
	availableEngines map[string]*PluginEngine
	runningEngines   map[string]*runningEngineContext
}

func NewPluginEngineRuntime(pluginDirs ...string) *pluginEngineRuntime {
	return &pluginEngineRuntime{
		pluginDirs:       pluginDirs,
		availableEngines: make(map[string]*PluginEngine),
		runningEngines:   make(map[string]*runningEngineContext),
	}
}

func (p *pluginEngineRuntime) LoadPlugins() error {
	// TODO: mutex

	plugins, err := LoadPlugins(p.pluginDirs...)
	if err != nil {
		return err
	}

	p.availableEngines = make(map[string]*PluginEngine)

	for _, pluginEngine := range plugins {
		name := pluginEngine.Engine.Name
		if _, has := p.availableEngines[name]; has == true {
			log.Warn("duplicate plugin %s", name)
		}
		p.availableEngines[name] = pluginEngine
	}

	return nil
}

func (p *pluginEngineRuntime) ListEngines() []*engines.EngineDescriptor {
	list := make([]*engines.EngineDescriptor, 0)
	for _, pluginEngine := range p.availableEngines {
		list = append(list, pluginEngine.Engine)
	}
	return list
}

func engineMatches(engine *PluginEngine, specification engines.Specification) bool {
	// TODO: implement a search
	return true
}

func (p *pluginEngineRuntime) getPluginEngine(engine *engines.EngineDescriptor) (*PluginEngine, bool) {
	pe, ok := p.availableEngines[engine.Name]
	return pe, ok
}

func (p *pluginEngineRuntime) FindEngines(specification engines.Specification) []*engines.EngineDescriptor {
	candidates := make([]*engines.EngineDescriptor, 0)

	for _, pluginEngine := range p.availableEngines {
		if engineMatches(pluginEngine, specification) {
			candidates = append(candidates, pluginEngine.Engine)
		}
	}

	return candidates
}

func (p *pluginEngineRuntime) FindEngineByName(name string) (*engines.EngineDescriptor, bool) {
	for _, pluginEngine := range p.availableEngines {
		e := pluginEngine.Engine
		if e.Name == name {
			return e, true
		}
	}
	return nil, false
}

func newRunningEngineContext() *runningEngineContext {
	ctx, cancelFunc := context.WithCancel(context.Background())
	started := new(sync.WaitGroup)
	started.Add(1)
	ctx = context.WithValue(ctx, "started", started)

	reCtx := &runningEngineContext{
		runningEngine: new(engines.RunningEngine),
		ctx:           ctx,
		cancel:        cancelFunc,
		started:       started,
	}

	return reCtx
}

func (p *pluginEngineRuntime) StartEngine(engine *engines.EngineDescriptor, spec messages.OperationSpec) (*engines.RunningEngine, error) {
	pluginEngine, ok := p.getPluginEngine(engine)
	if !ok {
		return nil, fmt.Errorf("could not find plugin engine for %s", engine.Name)
	}

	ctx := newRunningEngineContext()
	runtimeId := util.RandomString(15)
	ctx.runningEngine.RuntimeId = runtimeId
	ctx.ctx = context.WithValue(ctx.ctx, "attributes", spec.Attributes)

	// TODO: mutex
	p.runningEngines[runtimeId] = ctx

	startupObserver := make(chan messages.EngineAddress, 1)

	go func() {
		// TODO: create and pass a specification
		err := pluginEngine.Runner.Run(ctx.ctx, startupObserver, spec)

		if err != nil {
			log.Error("error running engine %s: %s", runtimeId, err)
			ctx.cancel()
		} else {
			log.Info("engine %s stopped", runtimeId)
		}
	}()

	log.Debug("waiting on engine %s(%s) to start", ctx.runningEngine.EngineDescriptor.Name, runtimeId)
	// wait for engine address to be determined by the plugin engine

	ctx.runningEngine.Address = <-startupObserver
	log.Debug("got engine address. waiting for started to be called")
	// FIXME: should also return on context cancellation
	ctx.started.Wait()
	log.Info("engine %s(%s) started", ctx.runningEngine.EngineDescriptor.Name, runtimeId)

	return ctx.runningEngine, nil
}

func (p *pluginEngineRuntime) StopEngine(engine *engines.RunningEngine) error {
	// TODO: mutex
	ctx, ok := p.runningEngines[engine.RuntimeId]
	if !ok {
		return fmt.Errorf("no such engine %s", engine.RuntimeId)
	}
	ctx.cancel()
	delete(p.runningEngines, engine.RuntimeId)
	return nil
}

func (p *pluginEngineRuntime) ListRunning() []*engines.RunningEngine {
	// TODO: mutex

	running := make([]*engines.RunningEngine, 0)
	for _, reCtx := range p.runningEngines {
		running = append(running, reCtx.runningEngine)
	}
	return running
}

func ParseSpec(specFilePath string) (*engines.EngineDescriptor, error) {
	data, err := ioutil.ReadFile(specFilePath)
	if err != nil {
		return nil, err
	}

	var engine engines.EngineDescriptor
	err = json.Unmarshal(data, &engine)
	if err != nil {
		return nil, err
	}

	return &engine, nil
}

func LoadPlugins(pluginDirs ...string) ([]*PluginEngine, error) {
	descriptors := make([]*engines.EngineDescriptor, 0)
	paths := make([]string, 0)
	plugins := make([]*PluginEngine, 0)

	// recursively collect engine descriptors
	collectFileDescriptors := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".spec.json") {
			return nil
		}

		descr, err := ParseSpec(path)
		if err != nil {
			log.Warn("error parsing stream spec of path %s: %s", path, err)
		}
		log.Info("found engine plugin descriptor: %s", descr)
		descriptors = append(descriptors, descr)
		paths = append(paths, path)
		return nil
	}
	for _, pluginDir := range pluginDirs {
		filepath.Walk(pluginDir, collectFileDescriptors)
	}

	// load plugin engines
	for i := 0; i < len(descriptors); i++ {
		descr := descriptors[i]
		descrFile := paths[i]

		if descr.Runtime == "" {
			if strings.HasSuffix(descrFile, ".so.spec.json") {
				descr.Runtime = "cogstream-go-plugin"
			} else {
				log.Error("cannot guess plugin type for descriptor at %s: %s", descrFile, descr)
				continue
			}
		}

		// engines using the go plugin system
		if descr.Runtime == "cogstream-go-plugin" {
			pluginPath := strings.TrimSuffix(descrFile, ".spec.json")
			engine, err := OpenPluginEngine(pluginPath)

			if err != nil {
				log.Warn("error loading plugin %s: %s", pluginPath, err)
				continue
			}
			plugins = append(plugins, engine)
			continue
		}

		// engines using the python cogstream.engine.srv
		if descr.Runtime == "cogstream-py" {
			engine, err := CreatePythonPluginEngine(descrFile, descr)
			if err != nil {
				log.Warn("error loading plugin %s: %s", descrFile, err)
				continue
			}
			plugins = append(plugins, engine)
			continue
		}

		log.Error("don't know how to handle plugin runtime %s of descriptor %s", descr.Runtime, descrFile)
	}

	return plugins, nil
}
