package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type responseBody struct {
	Id          int    `json:"id"`
	CleanedBody string `json:"cleaned_body"`
}

var counter int

func createChirp(w http.ResponseWriter, r *http.Request) {

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

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedBody := validateProfane(params.Body)

	counter++
	newId := counter

	respondWithJSON(w, 200, responseBody{Id: newId, CleanedBody: cleanedBody})

	// ..
	// params is a struct with data populated successfully
}

func validateProfane(s string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(s, " ")

	for i, word := range words {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(word) == strings.ToLower(profaneWord) {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}
