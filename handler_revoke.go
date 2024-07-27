package main

import (
	"fmt"
	"net/http"
	"strings"
)

type RevokeResponse struct {
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	fmt.Println("Received refresh token to be revoked is:", tokenString)

	status, err := cfg.DB.RevokeRefreshToken(tokenString)

	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		return
	}

	if status {
		respondWithJSON(w, 204, RevokeResponse{})
		return
	}

	respondWithError(w, http.StatusNotFound, "No token to be deleted")

}
