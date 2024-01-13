package mock

import (
	"errors"
	"lsat/macaroon"
	. "lsat/secrets"
)

type TestStore struct {
	users  map[UserId]Secret
	tokens map[macaroon.TokenID]macaroon.Token
}

func NewTestStore() TestStore {
	return TestStore{
		users:  make(map[UserId]Secret),
		tokens: make(map[macaroon.TokenID]macaroon.Token),
	}
}

func (store *TestStore) CreateUser() UserId {
	user := NewUserId()
	secret := NewSecret()

	store.users[user] = secret

	return user
}

func (store *TestStore) GetSecret(uid UserId) (Secret, error) {
	secret, ok := store.users[uid]

	if !ok {
		return Secret{}, errors.New("user not found!")
	}

	return secret, nil
}

func (store *TestStore) NewSecret(uid UserId) (Secret, error) {
	secret := NewSecret()

	store.users[uid] = secret

	return secret, nil
}

func (store *TestStore) StoreToken(id macaroon.TokenID, token macaroon.Token) error {
	store.tokens[id] = token
	return nil
}

func (store *TestStore) Tokens() *map[macaroon.TokenID]macaroon.Token {
	return &store.tokens
}
