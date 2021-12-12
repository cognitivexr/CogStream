package main

import (
	"cognitivexr.at/cogstream/mediator/cmd/mediator/app"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"cognitivexr.at/cogstream/mediator/pkg/mediator"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

type AbsPathList []string

func (l *AbsPathList) String() string {
	return strings.Join(*l, ",")
}

func (l *AbsPathList) Set(value string) (err error) {
	value, err = filepath.Abs(value)
	if err != nil {
		return
	}
	*l = append(*l, value)
	return
}

func main() {
	var pluginDirs AbsPathList

	hostPtr := flag.String("host", "0.0.0.0", "host to bind to")
	portPtr := flag.Int("port", 8191, "the server port")
	flag.Var(&pluginDirs, "plugins", "a directory containing engine plugins (can occur multiple times)")
	flag.Parse()

	if len(pluginDirs) == 0 {
		pluginDirs = append(pluginDirs, "engines/")
	}

	platform, err := mediator.NewPluginPlatform(pluginDirs...)
	if err != nil {
		log.Fatalf("error loading plugins from path: %s\n", pluginDirs, err)
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
