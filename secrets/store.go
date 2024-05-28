package secrets

import (
	"crypto/hmac"
	"crypto/sha256"
)

// SecretStore defines methods for managing secrets and tokens in a storage system.
type SecretStore interface {
	// NewSecret generates and returns a new secret associated with the provided user ID.
	// It returns an error if the operation fails.
	NewSecret(uid UserId) (Secret, error)

	// GetSecret retrieves the secret associated with the provided user ID.
	// It returns an error if the operation fails.
	GetSecret(uid UserId) (Secret, error)
}

type SecretFactory struct {
	root Secret
}

func NewSecretFactory() SecretFactory {
	return SecretFactory{
		root: NewSecret(),
	}
}

func NewStoreFromSecret(secret Secret) SecretFactory {
	return SecretFactory{
		root: secret,
	}
}

// Creates a new user ID.
func (store *SecretFactory) NewUser() UserId {
	return NewUserId()
}

func (store *SecretFactory) GetRoot() Secret {
	return store.root
}

func (store *SecretFactory) GetSecret(uid UserId) (Secret, error) {
	root := hmac.New(sha256.New, store.root[:])

	_, err := root.Write(uid[:])

	if err != nil {
		return Secret{}, err
	}

	return Secret(root.Sum(nil)), nil
}

func (store *SecretFactory) NewSecret(uid UserId) (Secret, error) {
	return store.GetSecret(uid)
}
