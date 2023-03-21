package mock

import (
	"encoding/binary"
	"errors"
	. "lsat/lightning"
	"math/rand"

	"github.com/lightningnetwork/lnd/lntypes"
)

func uint2bytes(n uint32) []byte {
	a := make([]byte, 32)
	binary.LittleEndian.PutUint32(a, n)
	return a
}

const Seed = 0

type TestChallenger struct {
	rand rand.Rand
}

func NewTestChallenger() TestChallenger {
	return TestChallenger{rand: *rand.New(rand.NewSource(Seed))}
}

func (node *TestChallenger) Challenge(price int64) (lntypes.Preimage, PaymentRequest, error) {
	preimage, err := lntypes.MakePreimage(uint2bytes(node.rand.Uint32()))

	if err != nil {
		return lntypes.Preimage{}, PaymentRequest{}, errors.New("failed to create preimage!")
	}

	rhash, _ := lntypes.MakeHash(uint2bytes(0))

	paymentRequest := PaymentRequest{
		RHash:       rhash,
		Invoice:     "",
		AddIndex:    0,
		PaymentAddr: make([]uint8, 0),
	}

	return preimage, paymentRequest, nil
}
