package mediator

import (
	"testing"
)

// TODO: extend tests

func Test_simpleHandshakeStore(t *testing.T) {
	store := NewSimpleHandshakeStore()

	hs := store.StartHandshake()
	println(hs.Created.Format("2006-01-02 15:04:05"))

	if hs1, ok := store.Get(hs.SessionId); ok {
		if hs1.SessionId != hs.SessionId {
			t.Error("ids do not match", hs.SessionId, hs1.SessionId)
		}
	} else {
		t.Error("no handshake with id", hs.SessionId)
	}
}
