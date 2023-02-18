package lightning

import (
	"time"

	"github.com/lightningnetwork/lnd/lntypes"
)

type PaymentRequest struct {
	Add_index    uint64
	Invoice      string
	R_hash       lntypes.Hash
	Payment_addr []uint8
}

type LightningNode interface {
	Pay(string) (lntypes.Preimage, error)
	CreateInvoice(int, time.Time, bool, string, lntypes.Preimage) (PaymentRequest, error)
	// SubscribeInvoice(r_hash lntypes.Hash) error
}

type Challenger interface {
	Challenge(int64) (lntypes.Preimage, PaymentRequest, error)
}

type ChallengeFactory struct {
	node LightningNode
}

func (node *ChallengeFactory) Challenge(price int64) (lntypes.Preimage, PaymentRequest, error) {
	return nil, PaymentRequest{}, nil
}

func (node *ChallengeFactory) Pay(invoice string) (lntypes.Preimage, error) {
	return lntypes.Preimage{}, nil
}

func (node *ChallengeFactory) CreateInvoice(value_msat int, expiry time.Time, private bool, memo string, preimage lntypes.Preimage) (PaymentRequest, error) {
	return PaymentRequest{}, nil
}
