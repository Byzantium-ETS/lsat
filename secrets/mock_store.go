package secrets

import (
	"errors"
	"lsat/auth"
)

type TestStore struct {
	users  map[UserId]Secret
	tokens map[auth.TokenID]auth.Token
}

func (store *TestStore) CreateUser() UserId {
	user := NewUserId()
	secret := NewSecret()

	store.users[user] = secret

	return user
}

func (store *TestStore) Secret(uid UserId) (Secret, error) {
	secret, ok := store.users[uid]

	if !ok {
		return Secret{}, errors.New("user not found!")
	}

	return secret, nil
}

func (store *TestStore) StoreToken(id auth.TokenID, token auth.Token) error {
	store.tokens[id] = token
	return nil
}

func (store *TestStore) Tokens() *map[auth.TokenID]auth.Token {
	return &store.tokens
}
