package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	fmt.Println("Received refresh token is:", tokenString)

	refreshToken, err := cfg.DB.GetRefreshToken(tokenString)

	if err != nil && refreshToken.Token == "" {
		respondWithError(w, http.StatusUnauthorized, "Token is not valid")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while fetching the refresh token")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy", IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)), Subject: fmt.Sprintf("%d", refreshToken.UserId)})

	signedToken, err := token.SignedString([]byte(cfg.jwtSecret))

	fmt.Println("new signed token:", signedToken)

	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Could not sign the token...")
		return
	}

	respondWithJSON(w, http.StatusOK, RefreshResponse{
		Token: signedToken,
	})

}
