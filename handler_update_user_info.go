package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UpdateResponse struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		fmt.Println("the token I'm returning is:", []byte(cfg.jwtSecret))
		return []byte(cfg.jwtSecret), nil
	})

	// an error will return if the token is not valid or expired
	if err != nil {
		log.Println("Something went wrong with the token:", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	userId, err := token.Claims.GetSubject()

	if err != nil {
		fmt.Println("Something went wrong while parsing the user id:", err)
		respondWithError(w, http.StatusInternalServerError, "Could not get the UserId...")
		return
	}

	fmt.Println("The user is:", userId)

	id, err := strconv.Atoi(userId)

	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Could not parse the user id...")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash the password")
		return
	}

	user, err := cfg.DB.UpdateUser(int(id), params.Email, string(hash))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user...")
		return
	}

	respondWithJSON(w, http.StatusOK, UpdateResponse{
		Id:       user.ID,
		Email:    user.Email,
		Password: string(hash),
	})
}
