package mediator

import (
	"cognitivexr.at/cogstream/pkg/api/messages"
	"math/rand"
	"time"
)

type Handshake struct {
	Created  time.Time
	Timeout  time.Duration
	Id       string
	Step     int
	Messages *messages.Messages
}

type HandshakeStore interface {
	SetTimeout(time.Duration)
	NewHandshake() *Handshake
	Get(id string) (*Handshake, bool)
}

type simpleHandshakeStore struct {
	storage map[string]*Handshake
	timeout time.Duration
}

func NewSimpleHandshakeStore() HandshakeStore {
	rand.Seed(time.Now().UnixNano())

	timeout, _ := time.ParseDuration("30s")

	store := &simpleHandshakeStore{
		storage: make(map[string]*Handshake),
		timeout: timeout,
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
		Step:     -1,
		Messages: messages.NewMessages(),
	}

	s.storage[id] = hs

	return hs
}

func (s *simpleHandshakeStore) Get(id string) (*Handshake, bool) {
	hs, ok := s.storage[id]

	return hs, ok
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
