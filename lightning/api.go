package lightning

import (
	"time"

	"github.com/lightningnetwork/lnd/lntypes"
)

type PaymentRequest struct {
	add_index    uint64
	invoice      string
	r_hash       lntypes.Hash
	payment_addr []uint8
}

type LightningNode interface {
	Pay(string) (lntypes.Preimage, error)
	CreateInvoice(int, time.Time, bool, string, lntypes.Preimage) (PaymentRequest, error)
	// SubscribeInvoice(r_hash lntypes.Hash) error
}

type Challenger interface {
	Challenge(lntypes.Preimage, int64) (PaymentRequest, error)
}

type Node struct {
	node LightningNode
}

func (node *Node) Challenge(preimage lntypes.Preimage, price int64) (PaymentRequest, error) {
	return PaymentRequest{}, nil
}

func (node *Node) Pay(invoice string) (lntypes.Preimage, error) {
	return lntypes.Preimage{}, nil
}

func (node *Node) CreateInvoice(value_msat int, expiry time.Time, private bool, memo string, preimage lntypes.Preimage) (PaymentRequest, error) {
	return PaymentRequest{}, nil
}
