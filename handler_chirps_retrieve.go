package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	queryAuthorId := r.URL.Query().Get("author_id")
	querySort := r.URL.Query().Get("sort")

	if queryAuthorId != "" {
		authorId, err := strconv.Atoi(queryAuthorId)

		if err != nil {
			fmt.Println("Something went wrong while converting", queryAuthorId, "into a number")
			return
		}
		dbChirps, err := cfg.DB.GetChirpsByAuthor(authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
			return
		}

		chirps := []Chirp{}
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:       dbChirp.ID,
				Body:     dbChirp.Body,
				AuthorId: dbChirp.AuthorId,
			})
		}

		if querySort == "asc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].ID < chirps[j].ID
			})
		}
		if querySort == "desc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].ID > chirps[j].ID
			})
		}

		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorId: dbChirp.AuthorId,
		})
	}

	if querySort == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	}
	if querySort == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsRetrieveById(w http.ResponseWriter, r *http.Request) {
	chirpId, err := strconv.Atoi(r.PathValue("id"))
	chirpId--

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch the ID")
		return
	}

	dbChirps, err := cfg.DB.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	if chirpId < 0 || chirpId >= len(dbChirps) {
		respondWithError(w, 404, "Chirp does not exist")
		return
	}

	chirp := dbChirps[chirpId]

	respondWithJSON(w, http.StatusOK, chirp)

}
