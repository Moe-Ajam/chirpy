package database

import (
	"errors"
	"fmt"
	"time"
)

type RefreshToken struct {
	Token      string    `json:"token"`
	ExpiryDate time.Time `json:"expiry_date"`
}

// returns the token if it exists, and an error if the something went wrong or the token doesn't exist
func (db *DB) GetRefreshToken(receivedRefreshToken string) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	for _, refreshToken := range dbStructure.RefreshTokens {
		if refreshToken.Token == receivedRefreshToken {
			return refreshToken, nil
		}
	}

	return RefreshToken{}, errors.New("Token doesn't exist")
}

func (db *DB) CreateRefreshToken(refreshTokenToBeSaved string) (RefreshToken, error) {
	dbStructure, err := db.loadDB()

	fmt.Println("loading db...")

	if err != nil {
		return RefreshToken{}, err
	}

	fmt.Println("Saving the token", refreshTokenToBeSaved)

	id := len(dbStructure.RefreshTokens) + 1

	fmt.Println("the id is:", id)

	refreshToken := RefreshToken{
		Token:      refreshTokenToBeSaved,
		ExpiryDate: time.Now().UTC().Add(60 * 24 * time.Hour),
	}

	fmt.Println("The refresh token to be saved:", refreshToken)

	dbStructure.RefreshTokens[id] = refreshToken

	err = db.writeDB(dbStructure)

	if err != nil {
		return RefreshToken{}, err
	}

	return refreshToken, nil
}
