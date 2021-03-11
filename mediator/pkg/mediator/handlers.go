package mediator

import (
	"cognitivexr.at/cogstream/pkg/api/messages"
	"fmt"
	"log"
)

//TODO add logic which establishes constraints
func RecordMessageHandler(hs *Handshake, msg messages.Message) (messages.Message, error) {
	defer incrementState(hs)
	log.Printf("handling message with code: %v", msg.GetCode())
	hs.InitiationCode = msg.GetCode()
	constraints := &messages.Constraints{}
	log.Printf("returning constraints: %v", constraints)
	return constraints, nil
}

//TODO add logic which establishes contents of agreement
func FormatMessageHandler(hs *Handshake, msg messages.Message) (messages.Message, error) {
	switch hs.InitiationCode {
	case messages.CodeRecord:
		defer incrementState(hs)
		agreement := &messages.RecordAgreement{
			Format: messages.Format{
				Resolution: "1x1",
			},
			URI: "a.b/c/d",
		}
		log.Printf("returning agreement: %v", agreement)
		return agreement, nil
	case messages.CodeExpose:
		defer incrementState(hs)
		agreement := &messages.ExposeAgreement{
			Format: messages.Format{
				Resolution: "",
			},
			URI: "",
		}
		log.Printf("returning agreement: %v", agreement)
		return agreement, nil
	case messages.CodeAnalyze:
		defer incrementState(hs)
		agreement := &messages.AnalyzeAgreement{
			Format: messages.Format{
				Resolution: "",
			},
		}
		log.Printf("returning agreement: %v", agreement)
		return agreement, nil
	default:
		return nil, fmt.Errorf("cannot establish handshake type")
	}
}

func incrementState(hs *Handshake) {
	hs.State = hs.State.Successor()
}
