package mock

import (
	"context"
	. "lsat/lightning"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
)

const Seed = 0

var node TestNode = TestNode{}
var challenger ChallengeFactory = NewChallenger(&node)

type TestNode struct{}

func (*TestNode) Pay(cx context.Context, invoice Invoice) (lntypes.Preimage, error) {
	return lntypes.Preimage{}, nil
}

func (*TestNode) CreateInvoice(cx context.Context, pr lnrpc.Invoice) (lnrpc.AddInvoiceResponse, error) {
	return lnrpc.AddInvoiceResponse{}, nil
}
