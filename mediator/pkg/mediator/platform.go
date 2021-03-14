package mediator

import (
	"cognitivexr.at/cogstream/api/messages"
)

type Platform interface {
	GetEngineFormatSpec(*HandshakeContext) (*messages.EngineFormatSpec, error)
	GetStreamSpec(*HandshakeContext) (*messages.StreamSpec, error)
}

type DummyPlatform struct {
}

func (d *DummyPlatform) GetEngineFormatSpec(_ *HandshakeContext) (*messages.EngineFormatSpec, error) {
	attributes := make(messages.Attributes)
	attributes.Set("foo", "bar")
	return &messages.EngineFormatSpec{Attributes: attributes}, nil
}

func (d *DummyPlatform) GetStreamSpec(hs *HandshakeContext) (*messages.StreamSpec, error) {
	attributes := make(messages.Attributes)
	attributes.Set("sessionId", hs.SessionId)
	attributes.Set("config.width", "640")
	attributes.Set("config.height", "360")
	attributes.Set("config.colormode", "1")
	return &messages.StreamSpec{EngineAddress: "127.0.0.1:53210", Attributes: attributes}, nil
}
