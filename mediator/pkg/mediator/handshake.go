package mediator

import (
	"cognitivexr.at/cogstream/pkg/api/messages"
	"fmt"
	"math/rand"
	"time"
)

type HandshakeState int

const (
	Empty HandshakeState = iota
	Operation
	Constraints
	Format
	Agreement
	Undefined
)

type Handshake struct {
	Created  time.Time
	Timeout  time.Duration
	Id       string
	State    HandshakeState
	Messages *messages.Messages
}

type HandshakeStore interface {
	SetTimeout(time.Duration)
	NewHandshake() *Handshake
	Get(id string) (*Handshake, bool)
	VerifyAndAddMessage(id string, message messages.Message) (HandshakeState, error)
}

type simpleHandshakeStore struct {
	storage      map[string]*Handshake
	timeout      time.Duration
	successorMap map[HandshakeState][]messages.MessageCode
}

func NewSimpleHandshakeStore() HandshakeStore {
	rand.Seed(time.Now().UnixNano())

	timeout, _ := time.ParseDuration("30s")

	store := &simpleHandshakeStore{
		storage: make(map[string]*Handshake),
		timeout: timeout,
		successorMap: map[HandshakeState][]messages.MessageCode{
			Empty:       {messages.CodeExpose, messages.CodeRecord, messages.CodeAnalyze},
			Operation:   {messages.CodeConstraints},
			Constraints: {messages.CodeFormat},
			Format:      {messages.CodeExposeAgreement, messages.CodeRecordAgreement, messages.CodeAnalyzeAgreement},
		},
	}

	// FIXME not ideal
	go func() {
		for range time.Tick(timeout) {
			store.expire()
		}
	}()

	return store
}

func (s *simpleHandshakeStore) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

func (s *simpleHandshakeStore) NewHandshake() *Handshake {
	id := s.nextHandshakeId()

	hs := &Handshake{
		Created:  time.Now(),
		Timeout:  s.timeout,
		Id:       id,
		State:    Empty,
		Messages: messages.NewMessages(),
	}

	s.storage[id] = hs

	return hs
}

func (s *simpleHandshakeStore) Get(id string) (*Handshake, bool) {
	hs, ok := s.storage[id]

	return hs, ok
}
func (s *simpleHandshakeStore) VerifyAndAddMessage(id string, message messages.Message) (HandshakeState, error) {
	hs, ok := s.storage[id]
	if !ok {
		return Undefined, fmt.Errorf("could not find handshake with id: %v", id)
	}

	if s.isAllowedMessageForState(hs.State, message) {
		hs.Messages.Add(message)
	} else {
		return Undefined, fmt.Errorf("message not allowed for state")
	}

	hs.State += 1

	return hs.State, nil
}

func (s *simpleHandshakeStore) isAllowedMessageForState(state HandshakeState, message messages.Message) bool {
	codes := s.successorMap[state]
	for _, code := range codes {
		if message.GetCode() == code {
			return true
		}
	}
	return false
}

func (s *simpleHandshakeStore) nextHandshakeId() string {
	// TODO: create random id
	return randomString(15)
}

func (s *simpleHandshakeStore) expire() {
	now := time.Now()
	threshold := now.Add(-s.timeout)

	for id, hs := range s.storage {
		if hs.Created.After(threshold) {
			delete(s.storage, id)
		}
	}
}

const runes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(l int) string {
	bytes := make([]byte, l)
	n := len(runes)
	for i := 0; i < l; i++ {
		bytes[i] = runes[rand.Intn(n)]
	}
	return string(bytes)
}
