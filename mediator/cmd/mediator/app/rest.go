package app

import (
	"cognitivexr.at/cogstream/api/messages"
	mediator "cognitivexr.at/cogstream/pkg/mediator"
	"encoding/json"
	"fmt"
	"net/http"
)

var Mediator *mediator.Mediator = nil

func InitEndpoints() {
	http.HandleFunc("/init/expose", initExpose)
	http.HandleFunc("/init/record", initRecord)
	http.HandleFunc("/init/analyze", initAnalyze)
	http.HandleFunc("/step", step)

	Mediator = mediator.NewMediator()
}

func initExpose(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{}")
	fmt.Println("Endpoint Hit: /init/expose")
}

func initRecord(w http.ResponseWriter, r *http.Request) {
	var recordMessage messages.Record

	err := json.NewDecoder(r.Body).Decode(&recordMessage)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hs := Mediator.StartHandshake()
	hs.Messages.Add(&recordMessage)

	fmt.Fprintf(w, "{'next': 'http://localhost:8191/step?session=%s'}", hs.Id)
	fmt.Printf("Endpoint Hit: /init/record: %v\n", recordMessage)
}

func initAnalyze(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{}")
	fmt.Println("Endpoint Hit: /init/analyze")
}

func step(w http.ResponseWriter, r *http.Request) {
	// next step in the negotiation
	session := r.URL.Query().Get("session")
	if session == "" {
		http.Error(w, "no session", http.StatusBadRequest)
		return
	}

	hs, ok := Mediator.Handshakes.Get(session)
	if !ok {
		http.Error(w, "no session", http.StatusBadRequest)
		return
	}

	// TODO: based on the current state of the handshake, determine what message type it should be
	println(hs.Created.Format("2006-01-03 15:04"))
}
