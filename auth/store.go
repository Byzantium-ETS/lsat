package auth

import (
	"encoding/json"
	"fmt"
	"lsat/macaroon"
	"os"
	"path/filepath"

	"github.com/lightningnetwork/lnd/lntypes"
)

// Constants
const (
	baseFileName = "l402.token."
)

type tokenJSON struct {
	Macaroon macaroon.MacaroonJSON `json:"macaroon"`
	Preimage string                `json:"preimage"`
}

// TokenStore defines the interface for storing and retrieving tokens.
type TokenStore interface {
	// Stores the provided token with the specified ID in the store.
	StoreToken(macaroon.TokenId, macaroon.Token) error

	// Returns a reference to the token stored in the store for the specified ID.
	GetToken(macaroon.TokenId) (*macaroon.Token, error)

	// Removes the token from the store
	RemoveToken(macaroon.TokenId) (*macaroon.Token, error)
}

// LocalStore implements the TokenStore interface using local file storage.
type LocalStore struct {
	directory string
}

// Create a new LocalStore.
func NewStore(directory string) (*LocalStore, error) {
	// If the target path for the token store doesn't exist, then we'll
	// create it now before we proceed.
	if !fileExists(directory) {
		if err := os.MkdirAll(directory, 0700); err != nil {
			return nil, err
		}
	}

	return &LocalStore{directory}, nil
}

// Saves the token to a file.
func (store *LocalStore) StoreToken(id macaroon.TokenId, token macaroon.Token) error {
	// Construct the file path
	filePath := store.FilePath(id)

	storedToken := tokenJSON{
		Macaroon: token.Macaroon.ToJSON(),
		Preimage: token.Preimage.String(),
	}

	// Marshal the token to JSON
	data, err := json.MarshalIndent(storedToken, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %v", err)
	}

	fmt.Println(string(data))

	// Write the JSON data to the file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write token to file: %v", err)
	}

	return nil
}

// GetToken reads the token from a file where it should be saved, unmarshals it, and returns the token object.
func (store *LocalStore) GetToken(id macaroon.TokenId) (*macaroon.Token, error) {
	// Construct the file path
	filePath := store.FilePath(id)

	return store.GetTokenFromPath(filePath)
}

// GetToken reads the token from a file where it should be saved, unmarshals it, and returns the token object.
func (store *LocalStore) GetTokenFromPath(filePath string) (*macaroon.Token, error) {
	// Read the JSON data from the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %v", err)
	}

	// Unmarshal the JSON data into a Token object
	var token tokenJSON
	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %v", err)
	}

	// Decode the preimage.
	preimage, err := lntypes.MakePreimageFromStr(token.Preimage)
	if err != nil {
		return nil, err
	}

	// Build a typed macaroon from the JSON object.
	mac, err := token.Macaroon.Unwrap()
	if err != nil {
		return nil, err
	}

	return &macaroon.Token{
		Macaroon: mac,
		Preimage: preimage,
	}, nil
}

func (store *LocalStore) RemoveToken(id macaroon.TokenId) (*macaroon.Token, error) {
	token, err := store.GetToken(id)
	if err != nil {
		return nil, err
	}

	err = os.Remove(store.FilePath(id))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (store *LocalStore) FilePath(id macaroon.TokenId) string {
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
