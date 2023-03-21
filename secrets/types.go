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
	arr := make([]byte, SecretSize)
	for i := 0; i < 4; i++ {
		n := rand.Uint64()
		binary.LittleEndian.PutUint64(arr[i*8:], n)
	}

	var secret Secret
	copy(secret[:32], arr[:32])

	return secret
}
