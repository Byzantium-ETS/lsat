package mock

import (
	"crypto/hmac"
	"crypto/sha256"
	"lsat/macaroon"
	. "lsat/secrets"
)

// TestStore represents a mock store for tokens and secrets.
type TestStore struct {
	root   Secret
	tokens map[macaroon.TokenID]macaroon.Token
}

func NewTestStore() TestStore {
	return TestStore{
		root:   NewSecret(),
		tokens: make(map[macaroon.TokenID]macaroon.Token),
	}
}

func NewTestStoreFromSecret(secret Secret) TestStore {
	return TestStore{
		root:   secret,
		tokens: make(map[macaroon.TokenID]macaroon.Token),
	}
}

// Creates a new user ID.
func (store *TestStore) CreateUser() UserId {
	return NewUserId()
}

func (store *TestStore) GetRoot() Secret {
	return store.root
}

func (store *TestStore) GetSecret(uid UserId) (Secret, error) {
	root := hmac.New(sha256.New, store.root[:])

	_, err := root.Write(uid[:])

	if err != nil {
		return Secret{}, err
	}

	return Secret(root.Sum(nil)), nil
}

func (store *TestStore) NewSecret(uid UserId) (Secret, error) {
	return store.GetSecret(uid)
}

func (store *TestStore) StoreToken(id macaroon.TokenID, token macaroon.Token) error {
	store.tokens[id] = token
	return nil
}

func (store *TestStore) Tokens() *map[macaroon.TokenID]macaroon.Token {
	return &store.tokens
}
