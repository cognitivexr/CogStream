package runtime

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/pkg/log"
	"cognitivexr.at/cogstream/pkg/util"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"plugin"
	"reflect"
	"strings"
	"sync"
)

type PluginEngine struct {
	Path   string
	Plugin *plugin.Plugin
	Runner engines.PluginEngineRunner
	Engine *engines.Engine
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

	return &PluginEngine{pluginPath, pluginObj, *runner, spec}, nil
}

type runningEngineContext struct {
	runningEngine *engines.RunningEngine
	ctx           context.Context
	cancel        context.CancelFunc
	started       *sync.WaitGroup
}

type pluginEngineRuntime struct {
	pluginDir        string
	availableEngines map[string]*PluginEngine
	runningEngines   map[string]*runningEngineContext
}

func NewPluginEngineRuntime(pluginDir string) *pluginEngineRuntime {
	return &pluginEngineRuntime{
		pluginDir:        pluginDir,
		availableEngines: make(map[string]*PluginEngine),
		runningEngines:   make(map[string]*runningEngineContext),
	}
}

func (p *pluginEngineRuntime) LoadPlugins() error {
	// TODO: mutex

	plugins, err := LoadPlugins(p.pluginDir)
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

func (p *pluginEngineRuntime) ListEngines() []*engines.Engine {
	list := make([]*engines.Engine, 0)
	for _, pluginEngine := range p.availableEngines {
		list = append(list, pluginEngine.Engine)
	}
	return list
}

func engineMatches(engine *PluginEngine, specification engines.Specification) bool {
	// TODO: implement a search
	return true
}

func (p *pluginEngineRuntime) getPluginEngine(engine *engines.Engine) (*PluginEngine, bool) {
	pe, ok := p.availableEngines[engine.Name]
	return pe, ok
}

func (p *pluginEngineRuntime) FindEngines(specification engines.Specification) []*engines.Engine {
	candidates := make([]*engines.Engine, 0)

	for _, pluginEngine := range p.availableEngines {
		if engineMatches(pluginEngine, specification) {
			candidates = append(candidates, pluginEngine.Engine)
		}
	}

	return candidates
}

func (p *pluginEngineRuntime) FindEngineByName(name string) (*engines.Engine, bool) {
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

func (p *pluginEngineRuntime) StartEngine(engine *engines.Engine) (*engines.RunningEngine, error) {
	pluginEngine, ok := p.getPluginEngine(engine)
	if !ok {
		return nil, fmt.Errorf("could not find plugin engine for %s", engine.Name)
	}

	ctx := newRunningEngineContext()
	runtimeId := util.RandomString(15)
	ctx.runningEngine.RuntimeId = runtimeId

	// TODO: how do we correctly determine the engine address?
	addr := "0.0.0.0:53210"
	ctx.runningEngine.Address = messages.EngineAddress(addr)

	// TODO: mutex
	p.runningEngines[runtimeId] = ctx

	go func() {
		// TODO: create and pass a specification
		err := pluginEngine.Runner.Run(ctx.ctx, engine, addr)

		if err != nil {
			log.Error("error running engine %s: %s", runtimeId, err)
			ctx.cancel()
		} else {
			log.Info("engine %s stopped", runtimeId)
		}
	}()

	log.Info("waiting on engine %s(%s) to start", ctx.runningEngine.Engine.Name, runtimeId)
	// FIXME: should also return on context cancellation
	ctx.started.Wait()
	log.Info("engine %s(%s) started", ctx.runningEngine.Engine.Name, runtimeId)

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

func ParseSpec(specFilePath string) (*engines.Engine, error) {
	data, err := ioutil.ReadFile(specFilePath)
	if err != nil {
		return nil, err
	}

	var engine engines.Engine
	err = json.Unmarshal(data, &engine)
	if err != nil {
		return nil, err
	}

	return &engine, nil
}

func LoadPlugins(pluginDir string) ([]*PluginEngine, error) {
	dir, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		return nil, err
	}

	plugins := make([]*PluginEngine, 0)

	for _, f := range dir {
		if f.IsDir() {
			continue
		}

		if !strings.HasSuffix(f.Name(), ".spec.json") {
			continue
		}

		pluginFile := strings.TrimSuffix(f.Name(), ".spec.json")
		pluginPath := path.Join(pluginDir, pluginFile)

		log.Info("loading plugin %s", pluginPath)
		engine, err := OpenPluginEngine(pluginPath)
		if err != nil {
			log.Warn("error loading plugin %s: %s", pluginPath, err)
			continue
		}
		plugins = append(plugins, engine)
	}

	return plugins, nil
}
