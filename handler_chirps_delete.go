package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpId)

	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while parsing the token")
		fmt.Println(err.Error())
		return
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while getting the user ID")
		return
	}
	authorId, err := strconv.Atoi(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		fmt.Println(err.Error())
		return
	}

	if authorId != chirp.AuthorId {
		respondWithError(w, 403, "Operation not allowed for this user")
		return
	}

	err = cfg.DB.DeleteChirp(chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while delting the chirp")
		fmt.Println(err.Error())
	}

	respondWithJSON(w, 204, "Chirp deleted!")

}
