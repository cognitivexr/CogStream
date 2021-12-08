package main

import (
	"cognitivexr.at/cogstream/mediator/cmd/mediator/app"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"cognitivexr.at/cogstream/mediator/pkg/mediator"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	hostPtr := flag.String("host", "0.0.0.0", "host to bind to")
	portPtr := flag.Int("port", 8191, "the server port")
	pluginDirPtr := flag.String("engine-dir", "engines/",
		"the directory containing engine plugins")

	flag.Parse()

	platform, err := mediator.NewPluginPlatform(*pluginDirPtr)
	if err != nil {
		log.Fatalf("could not load plugins from %s: %s\n", *pluginDirPtr, err)
		return
	}

	// for debugging purposes
	engines, err := platform.ListAvailableEngines()
	doc, err := mediator.MarshalAvailableEngines(engines)
	log.Info("plugins found: %s", doc)

	wsm := app.NewWebsocketMediator(mediator.NewSimpleHandshakeStore(), platform)
	wsm.AddOperationRequestHandler(mediator.DefaultOperationHandler)
	wsm.AddFormatEstablishmentHandler(mediator.DefaultFormatHandler)
	http.Handle("/", wsm)

	addr := fmt.Sprintf("%s:%d", *hostPtr, *portPtr)

	log.Info("starting mediator on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
