package main

import "sync"

type DB struct {
	path string
	mux  sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]ResponseBody `json:"chirps"`
}

// Creates a new database connection
// Creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {

}

func (db *DB) CreateChirp(body string) (Chirp, error) {

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {

}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {

}
