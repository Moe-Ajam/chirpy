package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
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
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters...")
		return
	}

	user, exists, err := cfg.DB.GetUserByEmail(params.Email)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy", IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)), Subject: fmt.Sprintf("%d", user.ID)})

	signedToken, err := token.SignedString([]byte(cfg.jwtSecret))

	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Could not sign the token...")
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

	randomBytes := make([]byte, 32)

	// turns the input into random values
	_, err = rand.Read(randomBytes)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate random bytes for the refresh token")
		return
	}

	encodedToken := hex.EncodeToString(randomBytes)

	_, err = cfg.DB.CreateRefreshToken(encodedToken, user.Email, user.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not save the refresh token to Database")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:           user.ID,
		Email:        user.Email,
		Token:        signedToken,
		RefreshToken: encodedToken,
		IsChirpyRed:  user.IsChirpyRed,
	})
}
