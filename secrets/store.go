package secrets

import (
	"lsat/auth"

	"github.com/lightningnetwork/lnd/lntypes"
)

type Secret = lntypes.Hash

type UserId = int32

type SecretStore interface {
	CreateSecret() (Secret, error)
	GetSecret(uid UserId) (Secret, error)          // S'il n'y a pas de RootKey pour l'utilisateur, il sera créé
	StoreToken(uid UserId, token auth.Token) error // Les tokens peuvent être conservé pour des raisons d'archivage
}
