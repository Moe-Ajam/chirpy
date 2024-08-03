package database

import (
	"errors"
	"fmt"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"passsword"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email string, passwordHash string) (User, error) {
	dbStructure, err := db.loadDB()

	fmt.Println("loading db...")

	if err != nil {
		return User{}, err
	}

	fmt.Println("authenticating for the email:", email, "and hashed password:", passwordHash)

	_, exists, err := db.GetUserByEmail(email)

	fmt.Println("User exists??", exists)

	if err != nil {
		return User{}, err
	}

	if exists {
		fmt.Println("User already exists")
		return User{}, nil
	}

	id := len(dbStructure.Users) + 1

	user := User{
		ID:          id,
		Email:       email,
		Password:    passwordHash,
		IsChirpyRed: false,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, false, err
	}

	// user, ok := dbStructure.Users[id]

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, true, nil
		}
	}

	fmt.Println("User with email:", email, "not found")

	return User{}, false, nil
}

func (db *DB) UpdateUser(id int, newEmail string, newPassword string) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	user := User{
		ID:       id,
		Email:    newEmail,
		Password: newPassword,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpgradeToChirpyRed(id int) error {
	dbStructure, err := db.loadDB()

	if err != nil {
		fmt.Println(err)
		return err
	}

	existingUser, exists := dbStructure.Users[id]

	if !exists {
		return errors.New("User does not exist")
	}
	existingUser.IsChirpyRed = true

	dbStructure.Users[id] = existingUser

	err = db.writeDB(dbStructure)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
