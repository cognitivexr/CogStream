package app

import (
	"cognitivexr.at/cogstream/pkg/mediator"
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
	//hs := mh.Mediator.StartHandshake()
	//var message messages.Message
	//
	//mt, messageReader, err := c.NextReader()
	//if err != nil {
	//	log.Printf("cannot read ws message: %v", err)
	//}
	//err = json.NewDecoder(messageReader).Decode(&message)
	//if err != nil {
	//	log.Printf("cannot decode message: %v", err)
	//}
	//log.Printf("got content: %v", message.Content)
	//
	//for {
	//	mt, rawMessage, err := c.ReadMessage()
	//	if err != nil {
	//		log.Println("cannot read ws message:", err)
	//		break
	//	}
	//	reader := bytes.NewReader(rawMessage)
	//	err = json.NewDecoder(reader).Decode(&envelope)
	//	if err != nil {
	//		log.Printf("cannot decode envelope: %v", err)
	//	}
	//	log.Printf("got data string: '%s'", string(envelope.Data))
	//	message, err := messages.UnmarshalJSONMessage(envelope.Type, bytes.NewReader(envelope.Data))
	//	if err != nil {
	//		log.Printf("cannot unmarshal message data: %v", err)
	//	}
	//	reply, err := mh.Mediator.ProcessMessage(hs.Id, message)
	//	if err != nil {
	//		log.Printf("cannot add message to handshake: %v", err)
	//	}
	//	w, err := c.NextWriter(mt)
	//	if err != nil {
	//		log.Printf("cannot open next writer: %v", err)
	//	}
	//	err = messages.WriteJSONMessage(reply, w)
	//	if err != nil {
	//		log.Printf("cannot send reply: %v", err)
	//		break
	//	}
	//}
}
