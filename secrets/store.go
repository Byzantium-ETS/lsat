package secrets

import (
	"encoding/binary"
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
