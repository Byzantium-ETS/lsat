package mock

import (
	"errors"
	"lsat/auth"
	"lsat/macaroon"
)

const (
	tokenErr = "could not find the token"
)

// LocalStore implements the TokenStore interface using local file storage.
type TestStore struct {
	tokens map[macaroon.TokenID]macaroon.Token
}

func NewStore() auth.TokenStore {
	return &TestStore{
		tokens: make(map[macaroon.TokenID]macaroon.Token),
	}
}

// Saves the token to a file at folderPath/baseFileName+id.Hash.
func (store *TestStore) StoreToken(id macaroon.TokenID, mac macaroon.Token) error {
	store.tokens[id] = mac
	return nil
}

// Reads the token from a file where it should be saved, unmarshals it, and returns the token object.
func (store *TestStore) GetToken(id macaroon.TokenID) (*macaroon.Token, error) {
	token, ok := store.tokens[id]

	if !ok {
		return nil, errors.New(tokenErr)
	}

	return &token, nil
}

// Remove the token from the store.
func (store *TestStore) RemoveToken(id macaroon.TokenID) (*macaroon.Token, error) {
	token, ok := store.tokens[id]

	if ok {
		delete(store.tokens, id)
	}

	if !ok {
		return nil, errors.New(tokenErr)
	}

	return &token, nil
}
