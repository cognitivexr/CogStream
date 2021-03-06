package mediator

import (
	"testing"
)

// TODO: extend tests

func Test_simpleHandshakeStore(t *testing.T) {
	store := NewSimpleHandshakeStore()

	hs := store.NewHandshake()
	println(hs.Created.Format("2006-01-02 15:04:05"))

	if hs1, ok := store.Get(hs.Id); ok {
		if hs1.Id != hs.Id {
			t.Error("ids do not match", hs.Id, hs1.Id)
		}
	} else {
		t.Error("no handshake with id", hs.Id)
	}
}
