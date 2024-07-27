package database

import (
	"errors"
	"fmt"
	"time"
)

type RefreshToken struct {
	Token      string    `json:"token"`
	Email      string    `json:"email"`
	UserId     int       `json:"user_id"`
	ExpiryDate time.Time `json:"expiry_date"`
}

// returns the token if it exists, and an error if the something went wrong or the token doesn't exist
func (db *DB) GetRefreshToken(receivedRefreshToken string) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	for _, refreshToken := range dbStructure.RefreshTokens {
		if refreshToken.Token == receivedRefreshToken && refreshToken.ExpiryDate.Sub(time.Now().UTC()) > 0 {
			return refreshToken, nil
		}
	}

	return RefreshToken{}, errors.New("Token doesn't exist")
}

func (db *DB) CreateRefreshToken(refreshTokenToBeSaved string, email string, userId int) (RefreshToken, error) {
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
		Email:      email,
		UserId:     userId,
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

func (db *DB) RevokeRefreshToken(receivedRefreshToken string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	for ky, refreshToken := range dbStructure.RefreshTokens {
		if refreshToken.Token == receivedRefreshToken {
			delete(dbStructure.RefreshTokens, ky)
			db.writeDB(dbStructure)
			return true, nil
		}
	}

	return false, errors.New("Token doesn't exist")
}
