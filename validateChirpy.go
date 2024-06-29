package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type responseBody struct {
	CleanedBody string `json:"cleaned_body"`
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
	}

	// respondWithJSON(w, 200, responseBody{Valid: true})

	cleanedBody := validateProfane(params.Body)

	respondWithJSON(w, 200, responseBody{CleanedBody: cleanedBody})
	// TODO: remove if not needed
	// do i really need this ??
	// respondWithJSON(w, 200, params)

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
