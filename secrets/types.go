package secrets

import (
	"math/rand"
)

const SecretSize = 32

type Secret = [SecretSize]byte

type UserId = uint32

func NewUserId() UserId {
	return rand.Uint32()
}

func NewSecret() Secret {
	var secret Secret

	rand.Read(secret[:])

	return secret
}
