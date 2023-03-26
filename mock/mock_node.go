package mock

import (
	. "lsat/lightning"
	"time"

	"github.com/lightningnetwork/lnd/lntypes"
)

const Seed = 0

var node TestNode = TestNode{}
var challenger ChallengeFactory = NewChallenger(&node)

type TestNode struct{}

func (Node *TestNode) Pay(invoice string) (lntypes.Preimage, error) {
	return lntypes.Preimage{}, nil
}

func (Node *TestNode) CreateInvoice(valueMsat uint64, expiry time.Time, private bool, memo string, preimage lntypes.Preimage) (PaymentRequest, error) {
	return PaymentRequest{}, nil
}
