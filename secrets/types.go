package secrets

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const SecretSize = 32

type Secret [SecretSize]byte

type UserId [SecretSize]byte

func NewUserId() UserId {
	var uid UserId

	rand.Read(uid[:])

	return uid
}

func MakeUserId(newUserId []byte) (UserId, error) {
	nhlen := len(newUserId)
	if nhlen != SecretSize {
		return UserId{}, fmt.Errorf("invalid user_id length of %v, want %v",
			nhlen, SecretSize)
	}

	var uid UserId
	copy(uid[:], newUserId)

	return uid, nil
}

func NewSecret() Secret {
	var secret Secret

	rand.Read(secret[:])

	return secret
}

func (u UserId) String() string {
	return hex.EncodeToString(u[:])
}

func (s Secret) String() string {
	return hex.EncodeToString(s[:])
}
