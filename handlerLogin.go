package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = 24 * 60 * 60 // 24 hours
	}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters...")
		return
	}

	fmt.Println("expiry time is set to:", params.ExpiresInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy", IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(params.ExpiresInSeconds) * time.Second)), Subject: params.Email})

	signedToken, err := token.SignedString([]byte(cfg.jwtSecret))

	fmt.Println("signed token:", signedToken)

	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Could not sign the token...")
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
		Token: signedToken,
	})
}
