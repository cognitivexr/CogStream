package main

import (
	"cognitivexr.at/cogstream/cmd/mediator/app"
	"cognitivexr.at/cogstream/pkg/mediator"
	"log"
	"net/http"
)

func main() {
	wsm := app.NewWebsocketMediator(mediator.NewSimpleHandshakeStore(), &mediator.DummyPlatform{})
	wsm.AddOperationRequestHandler(mediator.DefaultOperationHandler)
	wsm.AddFormatEstablishmentHandler(mediator.DefaultFormatHandler)
	http.Handle("/", wsm)
	log.Fatal(http.ListenAndServe(":8191", nil))
}
