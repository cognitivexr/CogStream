package main

import (
	"cognitivexr.at/cogstream/cmd/mediator/app"
	"cognitivexr.at/cogstream/pkg/log"
	"cognitivexr.at/cogstream/pkg/mediator"
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

	wsm := app.NewWebsocketMediator(mediator.NewSimpleHandshakeStore(), platform)
	wsm.AddOperationRequestHandler(mediator.DefaultOperationHandler)
	wsm.AddFormatEstablishmentHandler(mediator.DefaultFormatHandler)
	http.Handle("/", wsm)

	addr := fmt.Sprintf("%s:%d", *hostPtr, *portPtr)

	log.Info("starting mediator on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
