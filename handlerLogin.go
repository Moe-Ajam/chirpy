package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}

	user, exists, err := cfg.DB.GetUserByEmail(params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not fetch user")
		return
	}

	if !exists {
		respondWithError(w, http.StatusUnauthorized, "Email does not exist")
		return
	}

	cmpResult := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))

	if cmpResult != nil {
		respondWithError(w, http.StatusUnauthorized, "Password is wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}
