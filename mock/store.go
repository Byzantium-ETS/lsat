package mock

import (
	"errors"
	"lsat/auth"
	"lsat/macaroon"
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

// StoreToken saves the token to a file at folderPath/baseFileName+id.Hash.
func (store *TestStore) StoreToken(id macaroon.TokenID, mac macaroon.Token) error {
	store.tokens[id] = mac
	return nil
}

// GetToken reads the token from a file where it should be saved, unmarshals it, and returns the token object.
func (store *TestStore) GetToken(id macaroon.TokenID) (*macaroon.Token, error) {
	token, ok := store.tokens[id]

	if ok == false {
		return nil, errors.New("could not find the token!")
	}

	return &token, nil
}
