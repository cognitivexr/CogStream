package rpc

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/mediator/pkg/runtime"
	"net"
	"net/http"
	"net/rpc"
)

const DefaultServiceName = "RemoteEngineRuntime"

type EngineRuntimeSkeleton struct {
	runtime.EngineRuntime
	service  string
	listener net.Listener
}

type EngineRuntimeStub struct {
	serverAddress string
	service       string
}

func NewEngineRuntimeStub(serverAddress string) *EngineRuntimeStub {
	return &EngineRuntimeStub{
		serverAddress: serverAddress,
		service:       DefaultServiceName,
	}
}

func NewEngineRuntimeSkeleton(delegate runtime.EngineRuntime) *EngineRuntimeSkeleton {
	return &EngineRuntimeSkeleton{delegate, DefaultServiceName, nil}
}

type StartEngineCommand struct {
	Engine *engines.EngineDescriptor
	Spec   messages.OperationSpec
}

type StartEngineResponse struct {
	RunningEngine *engines.RunningEngine
}

type StopEngineCommand struct {
	Engine *engines.RunningEngine
}

type StopEngineResponse struct {
}

func (e *EngineRuntimeSkeleton) Serve(address string) error {
	rpc.RegisterName(e.service, e)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	e.listener = listener
	return http.Serve(listener, nil)
}

func (e *EngineRuntimeSkeleton) Close() error {
	if e.listener == nil {
		return nil
	}
	return e.listener.Close()
}

func (e *EngineRuntimeSkeleton) InvokeStartEngine(cmd *StartEngineCommand, resp *StartEngineResponse) error {
	engine, err := e.StartEngine(cmd.Engine, cmd.Spec)
	resp.RunningEngine = engine
	return err
}

func (e *EngineRuntimeSkeleton) InvokeStopEngine(cmd *StopEngineCommand, _ *StopEngineResponse) error {
	return e.StopEngine(cmd.Engine)
}

func (s *EngineRuntimeStub) client() (*rpc.Client, error) {
	return rpc.DialHTTP("tcp", s.serverAddress)
}

func (s *EngineRuntimeStub) StartEngine(engine *engines.EngineDescriptor, spec messages.OperationSpec) (*engines.RunningEngine, error) {
	client, err := s.client()
	if err != nil {
		return nil, err
	}

	cmd := &StartEngineCommand{
		Engine: engine,
		Spec:   spec,
	}
	resp := &StartEngineResponse{}

	err = client.Call(s.service+".InvokeStartEngine", cmd, resp)
	if err != nil {
		return nil, err
	}

	return resp.RunningEngine, nil
}

func (s *EngineRuntimeStub) StopEngine(engine *engines.RunningEngine) error {
	client, err := s.client()
	if err != nil {
		return err
	}

	cmd := &StopEngineCommand{
		Engine: engine,
	}
	resp := &StopEngineResponse{}

	err = client.Call(s.service+".InvokeStopEngine", cmd, resp)
	return err
}

func (s *EngineRuntimeStub) ListRunning() []*engines.RunningEngine {
	panic("implement me")
}
