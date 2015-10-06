package main

import (
	"encoding/json"
	"healthcheck/healthcheck"
	"net/http"
)

type PollJson struct {
	Key string `json:"key"`
}

func httpCreatePoll(w http.ResponseWriter, r *http.Request) {
	key, _ := healthcheck.CreatePoll()
	poll := PollJson{Key: key}

	w.WriteHeader(http.StatusCreated)
	je := json.NewEncoder(w)
	je.Encode(&poll)
}

func httpShowPoll(w http.ResponseWriter, r *http.Request) {
}
