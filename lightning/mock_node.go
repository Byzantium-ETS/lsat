package lightning

import (
	"errors"
	"lsat/secrets"
	"math/rand"

	"github.com/lightningnetwork/lnd/lntypes"
)

const Seed = 0

type TestChallenger struct {
	rand rand.Rand
}

func NewTestChallenger() TestChallenger {
	return TestChallenger{rand: *rand.New(rand.NewSource(Seed))}
}

func (node *TestChallenger) Challenge(price int64) (lntypes.Preimage, PaymentRequest, error) {
	preimage, err := lntypes.MakePreimage(secrets.Uint2bytes(node.rand.Uint32()))

	if err != nil {
		return nil, PaymentRequest{}, errors.New("failed to create preimage!")
	}

	rhash, _ := lntypes.MakeHash(secrets.Uint2bytes(0))

	payment_request := PaymentRequest{
		R_hash:       rhash,
		Invoice:      "",
		Add_index:    0,
		Payment_addr: make([]uint8, 0),
	}

	return preimage, payment_request, nil
}
