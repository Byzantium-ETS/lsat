package mock

import (
	. "lsat/challenge"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
)

const Seed = 0

type TestChallenger struct{}

var challenger = TestChallenger{}

func (*TestChallenger) Challenge(price uint64) (ChallengeResult, error) {
	preimage := lntypes.Preimage(secrets.NewSecret())
	return ChallengeResult{Preimage: preimage, Invoice: lnrpc.AddInvoiceResponse{}}, nil
}
