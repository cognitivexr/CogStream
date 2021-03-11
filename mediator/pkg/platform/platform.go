package platform

import (
	"cognitivexr.at/cogstream/pkg/api/messages"
	"cognitivexr.at/cogstream/pkg/mediator"
)

type Platform interface {
	GetEngineFormatSpec(*mediator.HandshakeContext) (*messages.EngineFormatSpec, error)
	GetStreamSpec(*mediator.HandshakeContext) (*messages.StreamSpec, error)
}

type DummyPlatform struct {
}

func (d *DummyPlatform) GetEngineFormatSpec(_ *mediator.HandshakeContext) (*messages.EngineFormatSpec, error) {
	attributes := make(messages.Attributes)
	attributes.Set("foo", "bar")
	return &messages.EngineFormatSpec{Attributes: attributes}, nil
}

func (d *DummyPlatform) GetStreamSpec(_ *mediator.HandshakeContext) (*messages.StreamSpec, error) {
	attributes := make(messages.Attributes)
	attributes.Set("foo", "baz")
	return &messages.StreamSpec{Attributes: attributes}, nil
}
