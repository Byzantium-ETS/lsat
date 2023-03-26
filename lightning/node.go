package lightning

import (
	"lsat/secrets"
	"time"

	"github.com/lightningnetwork/lnd/lntypes"
)

type PaymentRequest struct {
	AddIndex    uint64
	Invoice     string
	RHash       lntypes.Hash
	PaymentAddr []uint8
}

type Node interface {
	Pay(invoice string) (lntypes.Preimage, error)                                                                                   // Pay a BOLT11 invoice.
	CreateInvoice(valueMsat uint64, expiry time.Time, private bool, memo string, preimage lntypes.Preimage) (PaymentRequest, error) // Create a BOLT11 invoice.
}

type Challenger interface {
	Challenge(price uint64) (lntypes.Preimage, PaymentRequest, error)
}

type ChallengeFactory struct {
	node Node
}

func NewChallenger(node Node) ChallengeFactory {
	return ChallengeFactory{node}
}

func (challenger *ChallengeFactory) Challenge(price uint64) (lntypes.Preimage, PaymentRequest, error) {
	secret := secrets.NewSecret()
	preimage, _ := lntypes.MakePreimage(secret[:])
	paymentRequest, err := challenger.node.CreateInvoice(price, time.Now().Add(time.Duration(time.Minute)), false, "", preimage)
	return preimage, paymentRequest, err
}

func (challenger *ChallengeFactory) Pay(invoice string) (lntypes.Preimage, error) {
	return challenger.node.Pay(invoice)
}

func (challenger *ChallengeFactory) CreateInvoice(valueMsat uint64, expiry time.Time, private bool, memo string, preimage lntypes.Preimage) (PaymentRequest, error) {
	return challenger.node.CreateInvoice(valueMsat, expiry, private, memo, preimage)
}
