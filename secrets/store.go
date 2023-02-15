package secrets

import (
	"lsat/auth"
)

type Secret = []byte

type UserId = int32

type SecretStore interface {
	CreateSecret() (Secret, error)
	GetSecret(uid UserId) (Secret, error)          // S'il n'y a pas de RootKey pour l'utilisateur, il sera créé
	StoreLSAT(uid UserId, token auth.LSAT) error // Les tokens peuvent être conservé pour des raisons d'archivage
}
