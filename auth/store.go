package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"lsat/macaroon"
	"os"
	"path/filepath"
)

// Constants
const (
	baseFileName = "l402.token."
)

// TokenStore defines the interface for storing and retrieving tokens.
type TokenStore interface {
	// StoreToken stores the provided token with the specified ID in the store.
	// It returns an error if the operation fails.
	StoreToken(macaroon.TokenID, macaroon.Token) error

	// GetToken returns a reference to the token stored in the store for the specified ID.
	GetToken(macaroon.TokenID) (*macaroon.Token, error)
}

// LocalStore implements the TokenStore interface using local file storage.
type LocalStore struct {
	directory string
}

// Create a new LocalStore.
func NewStore(directory string) (LocalStore, error) {
	// If the target path for the token store doesn't exist, then we'll
	// create it now before we proceed.
	if !fileExists(directory) {
		if err := os.MkdirAll(directory, 0700); err != nil {
			return LocalStore{}, err
		}
	}

	return LocalStore{directory}, nil
}

// StoreToken saves the token to a file at folderPath/baseFileName+id.Hash.
func (store *LocalStore) StoreToken(id macaroon.TokenID, mac macaroon.Token) error {
	// Construct the file path
	filePath := store.FilePath(id)

	// Marshal the token to JSON
	data, err := json.Marshal(mac)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %v", err)
	}

	// Write the JSON data to the file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write token to file: %v", err)
	}

	return nil
}

// GetToken reads the token from a file where it should be saved, unmarshals it, and returns the token object.
func (store *LocalStore) GetToken(id macaroon.TokenID) (*macaroon.Token, error) {
	// Construct the file path
	filePath := store.FilePath(id)

	// Read the JSON data from the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		err := fmt.Sprintf("failed to read token file: %v\n", err)
		return nil, errors.New(err)
	}

	// Unmarshal the JSON data into a Token object
	var mac macaroon.Token
	err = json.Unmarshal(data, &mac)
	if err != nil {
		err := fmt.Sprintf("failed to unmarshal token: %v\n", err)
		return nil, errors.New(err)
	}

	return &mac, nil
}

func (store *LocalStore) FilePath(id macaroon.TokenID) string {
	return filepath.Join(store.directory, baseFileName+id.Hash.String())
}

// fileExists returns true if the file exists, and false otherwise.
func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
