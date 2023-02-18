package secrets

import (
	"encoding/binary"
	"lsat/auth"
	"math/rand"
)

const SecretSize = 32

type Secret = [SecretSize]byte

type UserId = uint32

func NewUserId() UserId {
	return rand.Uint32()
}

func NewSecret() Secret {
	n := rand.Uint64()
	arr := make([]byte, SecretSize)
	binary.LittleEndian.PutUint64(arr, n)

	var secret Secret
	copy(arr[:], secret[:])

	return secret
}

func Uint2bytes(n uint32) []byte {
	a := make([]byte, 32)
	binary.LittleEndian.PutUint32(a, n)
	return a
}

type SecretStore interface {
	// NewSecret() (Secret, error)
	Secret(uid UserId) (Secret, error)                  // S'il n'y a pas de RootKey pour l'utilisateur, il sera créé
	StoreToken(id auth.TokenID, token auth.Token) error // Les tokens peuvent être conservé pour des raisons d'archivage
	Tokens() *map[auth.TokenID]auth.Token               // Tous les tokens d'un utilisateur
}
