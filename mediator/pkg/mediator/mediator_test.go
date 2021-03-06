package mediator

import (
	"cognitivexr.at/cogstream/pkg/api/messages"
	"testing"
)

func TestMediator_ProcessMessage(t *testing.T) {
	m := NewMediator()
	m.Handle(messages.CodeRecord, m.defaultReply)
	m.Handle(messages.CodeFormat, m.defaultReply)

	hs := m.StartHandshake()
	println(hs.Created.Format("2006-01-02 15:04:05"))

	_, err := m.ProcessMessage(hs.Id, &messages.Record{Format: "rgb", Name: "test"})
	if err != nil {
		t.Errorf("couldn't process 'record' message: %v", err)
	}
	if hs.State != Constraints {
		t.Error("state not 'Constraints' after sending record")
	}

	_, err = m.ProcessMessage(hs.Id, &messages.Format{Resolution: "128x128"})
	if err != nil {
		t.Errorf("couldn't process 'format' message: %v", err)
	}
	if hs.State != Agreement {
		t.Error("state not 'Agreement' after sending format")
	}
}
