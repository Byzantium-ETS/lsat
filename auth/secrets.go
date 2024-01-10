package auth

import (
	"lsat/macaroon"
	"lsat/secrets"
)

// SecretStore defines methods for managing secrets and tokens in a storage system.
type SecretStore interface {
	// NewSecret generates and returns a new secret associated with the provided user ID.
	// It returns an error if the operation fails.
	NewSecret(uid secrets.UserId) (secrets.Secret, error)

	// GetSecret retrieves the secret associated with the provided user ID.
	// It returns an error if the operation fails.
	GetSecret(uid secrets.UserId) (secrets.Secret, error)

	// StoreToken stores the provided token with the specified ID in the store.
	// It returns an error if the operation fails.
	StoreToken(id macaroon.TokenID, token macaroon.Token) error

	// Tokens returns a reference to the map of tokens stored in the store.
	// Note: The returned map should not be modified directly to ensure data consistency.
	Tokens() *map[macaroon.TokenID]macaroon.Token
}
