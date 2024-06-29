package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type responseBody struct {
	Valid bool `json:"valid"`
}

func validateChirpyRequestHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	fmt.Println("Validation request received...")

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	} else {
		respondWithJSON(w, 200, responseBody{Valid: true})
		return
	}

	// TODO: remove if not needed
	// do i really need this ??
	// respondWithJSON(w, 200, params)

	// ..
	// params is a struct with data populated successfully
}