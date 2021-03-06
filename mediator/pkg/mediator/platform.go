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

func (d *DummyPlatform) GetStreamSpec(_ *HandshakeContext) (*messages.StreamSpec, error) {
	attributes := make(messages.Attributes)
	attributes.Set("foo", "baz")
	return &messages.StreamSpec{Attributes: attributes}, nil
}
