package main

import (
	"encoding/json"
	"net/http"
)

func validateChirpyRequestHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
}
