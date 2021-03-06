package mediator

import (
	"cognitivexr.at/cogstream/api/messages"
	"math/rand"
	"time"
)

// TODO: which fields to hide?
// TODO: turn into real "Context"? isn't really used in anything concurrently so far
type HandshakeContext struct {
	Created          time.Time
	Timeout          time.Duration
	SessionId        string
	OperationSpec    *messages.OperationSpec
	EngineFormatSpec *messages.EngineFormatSpec
	ClientFormatSpec *messages.ClientFormatSpec
	StreamSpec       *messages.StreamSpec
	Alert            *messages.Alert
	Ok               bool
}

type HandshakeStore interface {
	SetTimeout(time.Duration)
	StartHandshake() *HandshakeContext
	Get(id string) (*HandshakeContext, bool)
}

type simpleHandshakeStore struct {
	storage map[string]*HandshakeContext
	timeout time.Duration
}

func NewSimpleHandshakeStore() HandshakeStore {
	//TODO: maybe replace with UUID?
	rand.Seed(time.Now().UnixNano())

	timeout, _ := time.ParseDuration("30m")

	store := &simpleHandshakeStore{
		storage: make(map[string]*HandshakeContext),
		timeout: timeout,
	}

	// FIXME not ideal
	//go func() {
	//	for range time.Tick(timeout) {
	//		store.expire()
	//	}
	//}()
	// FIXED:
	time.AfterFunc(timeout, func() {
		store.expire()
	})

	return store
}

func (s *simpleHandshakeStore) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

func (s *simpleHandshakeStore) StartHandshake() *HandshakeContext {
	id := s.nextSessionId()

	hs := &HandshakeContext{
		Created:   time.Now(),
		Timeout:   s.timeout,
		SessionId: id,
		Ok:        true,
	}

	s.storage[id] = hs

	return hs
}

func (s *simpleHandshakeStore) Get(id string) (*HandshakeContext, bool) {
	hs, ok := s.storage[id]

	return hs, ok
}

func (s *simpleHandshakeStore) nextSessionId() string {
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
