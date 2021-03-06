package app

import (
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/pkg/mediator"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

type WebsocketMediator struct {
	*mediator.Mediator
	upgrader websocket.Upgrader
}

func NewWebsocketMediator(store mediator.HandshakeStore, platform mediator.Platform) *WebsocketMediator {
	return &WebsocketMediator{
		mediator.NewMediator(store, platform),
		websocket.Upgrader{},
	}
}

func (wm *WebsocketMediator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		err := c.Close()
		log.Printf("while closing websocket connection: %v", err)
	}()

	hs := wm.StartHandshake()
	wm.startCommunication(hs, c)
}

//TODO alerts and disconnects
func (wm *WebsocketMediator) startCommunication(hs *mediator.HandshakeContext, c *websocket.Conn) {
	mt, decodedMsg, err := wm.nextMessage(c)
	if err != nil {
		log.Printf("cannot get next message from client: %v", err)
		return
	}
	opSpec, err := mediator.UnmarshalOperationSpec(decodedMsg)
	if err != nil {
		log.Printf("cannot unmarshal operation spec: %v", err)
		return
	}
	err = wm.RequestOperation(hs, &opSpec)
	if err != nil {
		log.Printf("cannot request operation: %v", err)
		return
	}
	reply, err := mediator.MarshalEngineFormatSpec(hs.EngineFormatSpec)
	if err != nil {
		log.Printf("cannot marshal engine format spec: %v", err)
		return
	}
	err = c.WriteMessage(mt, reply)
	if err != nil {
		log.Printf("cannot write reply: %v", err)
		return
	}

	mt, decodedMsg, err = wm.nextMessage(c)
	if err != nil {
		log.Printf("cannot get next message from client: %v", err)
		return
	}
	clientSpec, err := mediator.UnmarshalClientFormatSpec(decodedMsg)
	if err != nil {
		log.Printf("cannot unmarshal client format spec: %v", err)
		return
	}
	err = wm.EstablishFormat(hs, &clientSpec)
	if err != nil {
		log.Printf("cannot establish format: %v", err)
		return
	}
	reply, err = mediator.MarshalStreamSpec(hs.StreamSpec)
	if err != nil {
		log.Printf("cannot marshal engine format spec: %v", err)
		return
	}
	err = c.WriteMessage(mt, reply)
	if err != nil {
		log.Printf("cannot write reply: %v", err)
		return
	}
}

func (wm *WebsocketMediator) nextMessage(c *websocket.Conn) (int, messages.Message, error) {
	mt, messageReader, err := c.NextReader()
	if err != nil {
		return 0, messages.Message{}, fmt.Errorf("cannot read from connection: %v", err)
	}

	decodedMsg, err := wm.decodeMessage(messageReader)
	if err != nil {
		return 0, messages.Message{}, fmt.Errorf("cannot decode received message: %v", err)
	}
	return mt, decodedMsg, nil
}

func (wm *WebsocketMediator) decodeMessage(r io.Reader) (messages.Message, error) {
	var message messages.Message
	err := json.NewDecoder(r).Decode(&message)
	if err != nil {
		return messages.Message{}, fmt.Errorf("cannot decode message from reader: %v", err)
	}
	return message, nil
}
