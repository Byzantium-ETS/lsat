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
	Pay(invoice string) (lntypes.Preimage, error)
	CreateInvoice(value_msat int, expiry time.Time, private bool, memo string, preimage lntypes.Preimage) (PaymentRequest, error)
	SubscribeInvoice(r_hash lntypes.Hash) error
}
