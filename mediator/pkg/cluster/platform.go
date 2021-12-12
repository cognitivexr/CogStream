package cluster

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
)

type RemotePlatform struct {
	leader *Leader
}

func (p *RemotePlatform) EngineWorkerMap() map[*engines.EngineDescriptor]*WorkerConnection {
	list := make(map[*engines.EngineDescriptor]*WorkerConnection)

	for _, worker := range p.leader.ListWorkers() {
		for _, engine := range worker.info.Engines {
			list[engine] = worker
		}
	}

	return list
}

func (p *RemotePlatform) ListEngines() []*engines.EngineDescriptor {
	list := make([]*engines.EngineDescriptor, 0)

	for _, worker := range p.leader.ListWorkers() {
		for _, engine := range worker.info.Engines {
			list = append(list, engine)
		}
	}

	return list
}

func (p *RemotePlatform) FindEngines(engines.Specification) []*engines.EngineDescriptor {
	candidates := make([]*engines.EngineDescriptor, 0)

	for _, engine := range p.ListEngines() {
		// TODO implement search
		candidates = append(candidates, engine)
	}

	return candidates
}

func (p *RemotePlatform) FindEngineByName(name string) (*engines.EngineDescriptor, bool) {
	// TODO: disambiguate between nodes that have the same engine
	for _, engine := range p.ListEngines() {
		e := engine
		if e.Name == name {
			return e, true
		}
	}
	return nil, false
}

func (p *RemotePlatform) StartEngine(engine *engines.EngineDescriptor, spec messages.OperationSpec) (*engines.RunningEngine, error) {
	// TODO implement me: load balancing/scheduling over RPC stubs of nodes
	panic("implement me")
}

func (p *RemotePlatform) StopEngine(engine *engines.RunningEngine) error {
	// TODO implement me: find node that's running the engine and call the RPC stub
	panic("implement me")
}
