package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UpdateResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerUpdateUserInfo(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters...")
		return
	}

	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	fmt.Println("token came as:", tokenString)

	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		fmt.Println("the token I'm returning is:", []byte(cfg.jwtSecret))
		return []byte(cfg.jwtSecret), nil
	})

	// an error will return if the token is not valid or expired
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is invalid")
		return
	}

	userId, err := token.Claims.GetSubject()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get the UserId...")
		return
	}

	fmt.Println("The user is:", userId)

	// user, _, err := cfg.DB.GetUserByEmail(params.Email)

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:    1,
		Email: "2",
	})
}
