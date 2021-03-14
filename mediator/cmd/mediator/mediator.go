package main

import (
	"cognitivexr.at/cogstream/cmd/mediator/app"
	"cognitivexr.at/cogstream/pkg/mediator"
	"log"
	"net/http"
)

func main() {
	pluginDir := "../../engines/build" // FIXME: add configurable path

	platform, err := mediator.NewPluginPlatform(pluginDir)
	if err != nil {
		log.Fatalf("could not load plugins from %s: %s\n", pluginDir, err)
		return
	}

	wsm := app.NewWebsocketMediator(mediator.NewSimpleHandshakeStore(), platform)
	wsm.AddOperationRequestHandler(mediator.DefaultOperationHandler)
	wsm.AddFormatEstablishmentHandler(mediator.DefaultFormatHandler)
	http.Handle("/", wsm)
	log.Fatal(http.ListenAndServe(":8191", nil))
}
