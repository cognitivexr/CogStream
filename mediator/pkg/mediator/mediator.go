package mediator

import (
	"cognitivexr.at/cogstream/pkg/api/messages"
	"cognitivexr.at/cogstream/pkg/platform"
	"fmt"
)

type HandshakeStepHandler func(*HandshakeContext, platform.Platform) error

type Mediator struct {
	handshakes                  HandshakeStore
	operationRequestHandlers    []HandshakeStepHandler
	formatEstablishmentHandlers []HandshakeStepHandler
	platform                    platform.Platform
	handlersLocked              bool
}

func NewMediator() *Mediator {
	m := &Mediator{
		handshakes:                  NewSimpleHandshakeStore(),
		operationRequestHandlers:    make([]HandshakeStepHandler, 0),
		formatEstablishmentHandlers: make([]HandshakeStepHandler, 0),
		handlersLocked:              false,
		platform:                    &platform.DummyPlatform{},
	}

	m.AddOperationRequestHandler(DefaultOperationHandler)
	m.AddFormatEstablishmentHandler(DefaultFormatHandler)

	return m
}

func (m *Mediator) StartHandshake() *HandshakeContext {
	if !m.handlersLocked {
		m.handlersLocked = true
	}
	return m.handshakes.StartHandshake()
}

func (m *Mediator) RequestOperation(sessionId string, spec *messages.OperationSpec) error {
	hs, ok := m.handshakes.Get(sessionId)
	if !ok {
		return fmt.Errorf("no session with id: %v", sessionId)
	}
	if !hs.Ok {
		return fmt.Errorf("state is not ok: %v", hs)
	}
	if hs.OperationSpec != nil || hs.EngineFormatSpec != nil || hs.ClientFormatSpec != nil || hs.StreamSpec != nil {
		hs.Ok = false
		return fmt.Errorf("state is not empty: %v", hs)
	}

	hs.OperationSpec = spec

	for _, handler := range m.operationRequestHandlers {
		err := handler(hs, m.platform)
		if err != nil {
			hs.Ok = false
			return fmt.Errorf("cannot handle operation spec: %v", err)
		}
	}
	if hs.EngineFormatSpec == nil {
		hs.Ok = false
		return fmt.Errorf("no engine format set")
	}

	return nil
}

func (m *Mediator) EstablishFormat(sessionId string, spec *messages.ClientFormatSpec) error {
	hs, ok := m.handshakes.Get(sessionId)
	if !ok {
		return fmt.Errorf("no session with id: %v", sessionId)
	}
	if !hs.Ok {
		return fmt.Errorf("state is not ok: %v", hs)
	}
	if hs.OperationSpec == nil || hs.EngineFormatSpec == nil || hs.ClientFormatSpec != nil || hs.StreamSpec != nil {
		hs.Ok = false
		return fmt.Errorf("state is not ready for format establishment: %v", hs)
	}

	hs.ClientFormatSpec = spec
	for _, handler := range m.formatEstablishmentHandlers {
		err := handler(hs, m.platform)
		if err != nil {
			hs.Ok = false
			return fmt.Errorf("cannot handle client format spec: %v", err)
		}
	}
	if hs.StreamSpec == nil {
		hs.Ok = false
		return fmt.Errorf("no stream spec set")
	}

	return nil
}

func (m *Mediator) AddOperationRequestHandler(handler HandshakeStepHandler) bool {
	if m.handlersLocked {
		return false
	}

	m.operationRequestHandlers = append(m.operationRequestHandlers, handler)

	return true
}
func (m *Mediator) AddFormatEstablishmentHandler(handler HandshakeStepHandler) bool {
	if m.handlersLocked {
		return false
	}

	m.formatEstablishmentHandlers = append(m.formatEstablishmentHandlers, handler)

	return true
}
