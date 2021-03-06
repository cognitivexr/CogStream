package mediator

import (
	"cognitivexr.at/cogstream/pkg/api/messages"
	"fmt"
	"sync"
)

type Mediator struct {
	Handshakes      HandshakeStore
	successors      map[HandshakeState][]messages.MessageCode
	messageHandlers map[messages.MessageCode]func(*Handshake, messages.Message) (messages.Message, error)
	handlerMutex    sync.Mutex
}

func NewMediator() *Mediator {
	m := &Mediator{
		Handshakes: NewSimpleHandshakeStore(),
		successors: map[HandshakeState][]messages.MessageCode{
			Empty:       {messages.CodeExpose, messages.CodeRecord, messages.CodeAnalyze},
			Operation:   {messages.CodeConstraints},
			Constraints: {messages.CodeFormat},
			Format:      {messages.CodeExposeAgreement, messages.CodeRecordAgreement, messages.CodeAnalyzeAgreement},
		},
		messageHandlers: make(map[messages.MessageCode]func(*Handshake, messages.Message) (messages.Message, error)),
	}
	m.Handle(messages.CodeRecord, m.defaultReply)
	m.Handle(messages.CodeFormat, m.defaultReply)
	return m
}

func (m *Mediator) StartHandshake() *Handshake {
	return m.Handshakes.NewHandshake()
}

func (m *Mediator) ProcessMessage(id string, message messages.Message) (messages.Message, error) {
	hs, ok := m.Handshakes.Get(id)
	if !ok {
		return nil, fmt.Errorf("could not find handshake with id: %v", id)
	}

	if !m.isAllowedMessageForState(hs.State, message) {
		return nil, fmt.Errorf("message not allowed for state")
	}
	hs.Messages.Add(message)
	hs.State = hs.State.Successor()

	m.handlerMutex.Lock()
	handler := m.messageHandlers[message.GetCode()]
	m.handlerMutex.Unlock()

	if handler == nil {
		return nil, fmt.Errorf("no handler for message with code %v", message.GetCode())
	}
	reply, err := handler(hs, message)
	if err != nil {
		return nil, fmt.Errorf("couldn't create reply to message: %v", err)
	}

	return reply, nil
}

func (m *Mediator) Handle(code messages.MessageCode, handler func(*Handshake, messages.Message) (messages.Message, error)) {
	m.handlerMutex.Lock()
	defer m.handlerMutex.Unlock()
	m.messageHandlers[code] = handler
}

func (m *Mediator) isAllowedMessageForState(state HandshakeState, message messages.Message) bool {
	codes := m.successors[state]
	for _, code := range codes {
		if message.GetCode() == code {
			return true
		}
	}
	return false
}

func (m *Mediator) defaultReply(hs *Handshake, _ messages.Message) (messages.Message, error) {
	hs.State = hs.State.Successor()
	return &messages.Record{}, nil
}
