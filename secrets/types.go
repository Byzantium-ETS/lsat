package secrets

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const SecretSize = 32

// Secret is a fixed-size byte array representing a secret value.
type Secret [SecretSize]byte

// UserID is a fixed-size byte array representing a user identifier.
type UserID [SecretSize]byte

// NewUserId generates a new unique UserId.
func NewUserId() UserID {
	var uid UserID
	rand.Read(uid[:]) // Fill the byte array with random data.
	return uid
}

// MakeUserId creates a UserId from the provided byte slice.
func MakeUserId(newUserId []byte) (UserID, error) {
	nhlen := len(newUserId)
	if nhlen != SecretSize {
		return UserID{}, fmt.Errorf("invalid user_id length of %v, want %v",
			nhlen, SecretSize)
	}

	var user_id UserID
	copy(user_id[:], newUserId) // Copy the provided slice into the UserId byte array.

	return user_id, nil
}

// MakeUserIdFromStr creates a UserId from a hex string.
func MakeUserIdFromStr(newuser string) (UserID, error) {
	if len(newuser) != SecretSize*2 {
		return UserID{}, fmt.Errorf("invalid user_id string length of %v, "+
			"want %v", len(newuser), SecretSize*2)
	}

	user_id, err := hex.DecodeString(newuser)
	if err != nil {
		return UserID{}, err
	}

	return MakeUserId(user_id)
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

func (u UserID) String() string {
	return hex.EncodeToString(u[:])
}

func (s Secret) String() string {
	return hex.EncodeToString(s[:])
}
