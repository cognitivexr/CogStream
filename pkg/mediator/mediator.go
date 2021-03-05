package mediator

type Mediator struct {
	Handshakes HandshakeStore
}

func NewMediator() *Mediator {
	return &Mediator{
		Handshakes: NewSimpleHandshakeStore(),
	}
}

func (m *Mediator) StartHandshake() *Handshake {
	return m.Handshakes.NewHandshake()
}
