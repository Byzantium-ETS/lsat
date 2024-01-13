package mock

import (
	"lsat/challenge"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
)

type TestChallenger struct{}

func NewChallenger() challenge.Challenger {
	return &TestChallenger{}
}

func (*TestChallenger) Challenge(price uint64) (challenge.ChallengeResult, error) {
	preimage := lntypes.Preimage(secrets.NewSecret())
	return challenge.ChallengeResult{Preimage: preimage, PaymentRequest: lnrpc.AddInvoiceResponse{}}, nil
}
