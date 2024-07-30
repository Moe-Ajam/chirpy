package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerUpgradeChirpyRed(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserId int `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters...")
		return
	}

	if params.Event != "user.upgraded" {
		fmt.Println("Event came as:", params.Event)
		respondWithJSON(w, 204, "Doesn't really matter")
		return
	}

	fmt.Println("UserID that will be fetched:", params.Data.UserId)
	err = cfg.DB.UpgradeToChirpyRed(params.Data.UserId)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, 404, "User doesn't exist or something went wrong")
		return
	}

	respondWithJSON(w, 204, "User updated Successfully!")
}
