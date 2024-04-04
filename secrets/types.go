package secrets

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const SecretSize = 32

// Secret is a fixed-size byte array representing a secret value.
type Secret [SecretSize]byte

// UserId is a fixed-size byte array representing a user identifier.
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

func MakeSecret(NewSecret []byte) (Secret, error) {
	nhlen := len(NewSecret)
	if nhlen != SecretSize {
		return Secret{}, fmt.Errorf("invalid secret length of %v, want %v",
			nhlen, SecretSize)
	}

	var secret UserId
	copy(secret[:], NewSecret)

	return Secret(secret), nil
}

func (u UserId) String() string {
	return hex.EncodeToString(u[:])
}

func (s Secret) String() string {
	return hex.EncodeToString(s[:])
}
