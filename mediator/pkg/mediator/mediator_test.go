package mediator

import (
	"cognitivexr.at/cogstream/api/messages"
	"testing"
)

// TODO: extend tests

func TestMediator_ProcessMessage(t *testing.T) {
	shs := NewSimpleHandshakeStore()
	dp := DummyPlatform{}
	m := NewMediator(shs, &dp)
	m.AddOperationRequestHandler(DefaultOperationHandler)
	m.AddFormatEstablishmentHandler(DefaultFormatHandler)

	hs := m.StartHandshake()
	println(hs.Created.Format("2006-01-02 15:04:05"))

	err := m.RequestOperation(hs, &messages.OperationSpec{
		Code: "record",
		Attributes: messages.Attributes{
			"foo": {"bar", "baz"},
		},
	})
	if err != nil {
		t.Errorf("cannot request operation: %v", err)
	}
	if hs.OperationSpec == nil {
		t.Error("operation spec was not set in handshake context")
	}
	if hs.EngineFormatSpec == nil {
		t.Error("engine format spec was not set in handshake context")
	}

	err = m.EstablishFormat(hs, &messages.ClientFormatSpec{
		Attributes: messages.Attributes{
			"la": {"le", "lu"},
		},
	})
	if err != nil {
		t.Errorf("cannot establish format: %v", err)
	}
	if hs.ClientFormatSpec == nil {
		t.Error("client format spec was not set in handshake context")
	}
	if hs.StreamSpec == nil {
		t.Error("stream spec was not set in handshake context")
	}
}
