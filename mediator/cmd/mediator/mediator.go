package main

import (
	"cognitivexr.at/cogstream/cmd/mediator/app"
	"log"
	"net/http"
)

func main() {
	app.InitEndpoints()
	log.Fatal(http.ListenAndServe(":8191", nil))
}
