package main

import (
	"cognitivexr.at/cogstream/cmd/mediator/app"
	"log"
	"net/http"
)

func main() {
	mh := app.NewMediatorHandler()
	http.Handle("/", mh)
	log.Fatal(http.ListenAndServe(":8191", nil))
}
