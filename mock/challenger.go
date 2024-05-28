package mock

import (
	"lsat/challenge"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lntypes"
)

type testChallenger struct{}

func NewChallenger() challenge.Challenger {
	return &testChallenger{}
}

// The invoice will be the preimage for testing purposes.
func (*testChallenger) Challenge(price uint64) (challenge.InvoiceResponse, error) {
	preimage := lntypes.Preimage(secrets.NewSecret())
	return challenge.InvoiceResponse{
		Invoice:     preimage.String(),
		Preimage:    preimage,
		PaymentHash: preimage.Hash(),
	}, nil
}
