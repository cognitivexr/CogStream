package mediator

import (
	"cognitivexr.at/cogstream/api/messages"
	"fmt"
)

type HandshakeStepHandler func(*HandshakeContext, Platform) error

type Mediator struct {
	handshakes                  HandshakeStore
	operationRequestHandlers    []HandshakeStepHandler
	formatEstablishmentHandlers []HandshakeStepHandler
	platform                    Platform
	handlersLocked              bool
}

func NewMediator(store HandshakeStore, platform Platform) *Mediator {
	m := &Mediator{
		handshakes:                  store,
		operationRequestHandlers:    make([]HandshakeStepHandler, 0),
		formatEstablishmentHandlers: make([]HandshakeStepHandler, 0),
		handlersLocked:              false,
		platform:                    platform,
	}
	return m
}

func (m *Mediator) StartHandshake() *HandshakeContext {
	if !m.handlersLocked {
		m.handlersLocked = true
	}
	return m.handshakes.StartHandshake()
}

// RequestOperation takes a HandshakeContext and a messages.OperationSpec and adds the messages.OperationSpec
// and a messages.AvailableEngines to the HandshakeContext
func (m *Mediator) RequestOperation(hs *HandshakeContext, spec *messages.OperationSpec) error {
	if !hs.Ok {
		return fmt.Errorf("state is not ok: %v", hs)
	}
	if hs.OperationSpec != nil || hs.AvailableEngines != nil || hs.ClientFormatSpec != nil || hs.StreamSpec != nil {
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
	if hs.AvailableEngines == nil {
		hs.Ok = false
		return fmt.Errorf("no engine format set")
	}

	return nil
}

// EstablishFormat takes a HandshakeContext and a messages.ClientFormatSpec and adds the messages.ClientFormatSpec
// and a messages.StreamSpec to the HandshakeContext
func (m *Mediator) EstablishFormat(hs *HandshakeContext, spec *messages.ClientFormatSpec) error {
	if !hs.Ok {
		return fmt.Errorf("state is not ok: %v", hs)
	}
	if hs.OperationSpec == nil || hs.AvailableEngines == nil || hs.ClientFormatSpec != nil || hs.StreamSpec != nil {
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
