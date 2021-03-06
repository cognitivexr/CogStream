package app

import (
	"bytes"
	"cognitivexr.at/cogstream/pkg/api/messages"
	"cognitivexr.at/cogstream/pkg/mediator"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type MediatorHandler struct {
	*mediator.Mediator
	upgrader websocket.Upgrader
}

func NewMediatorHandler() *MediatorHandler {
	return &MediatorHandler{
		mediator.NewMediator(),
		websocket.Upgrader{},
	}
}

func (mh *MediatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := mh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()
	hs := mh.Mediator.StartHandshake()
	var envelope messages.Envelope
	for {
		mt, rawMessage, err := c.ReadMessage()
		if err != nil {
			log.Println("couldn't read ws message:", err)
			break
		}
		reader := bytes.NewReader(rawMessage)
		err = json.NewDecoder(reader).Decode(&envelope)
		if err != nil {
			log.Printf("couldn't decode envelope: %v", err)
		}
		log.Printf("got data string: '%s'", string(envelope.Data))
		message, err := messages.UnmarshalJSONMessage(envelope.Type, bytes.NewReader(envelope.Data))
		if err != nil {
			log.Printf("couldn't unmarshal message data: %v", err)
		}
		_, err = mh.Mediator.ProcessMessage(hs.Id, message)
		if err != nil {
			log.Printf("couldn't add message to handshake: %v", err)
		}
		err = c.WriteMessage(mt, []byte("not implemented"))
		if err != nil {
			log.Printf("couldn't send reply: %v", err)
			break
		}
	}
}
