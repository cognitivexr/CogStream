package mediator

import (
	"cognitivexr.at/cogstream/pkg/api/messages"
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

func TestRecordHandshake(t *testing.T) {
	store := NewSimpleHandshakeStore()

	hs := store.NewHandshake()
	println(hs.Created.Format("2006-01-02 15:04:05"))

	state, err := store.VerifyAndAddMessage(hs.Id, &messages.Record{Format: "rgb", Name: "test"})
	if err != nil {
		t.Errorf("couldn't add 'record' message: %v", err)
	}
	if state != Operation || hs.State != Operation {
		t.Error("state not 'Operation' after sending record")
	}
	state, err = store.VerifyAndAddMessage(hs.Id, &messages.Constraints{})
	if err != nil {
		t.Errorf("couldn't add 'constraint' message: %v", err)
	}
	if state != Constraints || hs.State != Constraints {
		t.Error("state not 'Constraints' after sending constraints")
	}
	state, err = store.VerifyAndAddMessage(hs.Id, &messages.Format{Resolution: "128x128"})
	if err != nil {
		t.Errorf("couldn't add 'format' message: %v", err)
	}
	if state != Format || hs.State != Format {
		t.Error("state not 'Format' after sending format")
	}
	state, err = store.VerifyAndAddMessage(hs.Id, &messages.RecordAgreement{Format: messages.Format{Resolution: "128x128"}, URI: randomString(8)})
	if err != nil {
		t.Errorf("couldn't add 'recordAgreement' message: %v", err)
	}
	if state != Agreement || hs.State != Agreement {
		t.Error("state not 'Agreement' after sending record agreement")
	}
}

func TestExposeHandshake(t *testing.T) {

}
