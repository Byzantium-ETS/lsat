package secrets

import (
	"crypto/rand"
)

const SecretSize = 32

type Secret = [SecretSize]byte

type UserId = [SecretSize]byte

func NewUserId() Secret {
	var id UserId

	rand.Read(id[:])

	return id
}

func NewSecret() Secret {
	var secret Secret

	rand.Read(secret[:])

	return secret
}
