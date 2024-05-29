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

// NewUserId generates a new unique UserId.
func NewUserId() UserId {
	var uid UserId
	rand.Read(uid[:]) // Fill the byte array with random data.
	return uid
}

// MakeUserId creates a UserId from the provided byte slice.
func MakeUserId(newUserId []byte) (UserId, error) {
	nhlen := len(newUserId)
	if nhlen != SecretSize {
		return UserId{}, fmt.Errorf("invalid user_id length of %v, want %v",
			nhlen, SecretSize)
	}

	var uid UserId
	copy(uid[:], newUserId) // Copy the provided slice into the UserId byte array.

	return uid, nil
}

// NewSecret generates a new random Secret.
func NewSecret() Secret {
	var secret Secret
	rand.Read(secret[:]) // Fill the byte array with random data.
	return secret
}

// MakeSecret creates a Secret from the provided byte slice.
func MakeSecret(newSecret []byte) (Secret, error) {
	nhlen := len(newSecret)
	if nhlen != SecretSize {
		return Secret{}, fmt.Errorf("invalid secret length of %v, want %v",
			nhlen, SecretSize)
	}

	var secret Secret
	copy(secret[:], newSecret) // Copy the provided slice into the Secret byte array.

	return secret, nil
}

func (u UserId) String() string {
	return hex.EncodeToString(u[:])
}

func (s Secret) String() string {
	return hex.EncodeToString(s[:])
}
