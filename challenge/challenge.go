package challenge

import (
	"context"
	"lsat/secrets"
	"time"

	"github.com/lightningnetwork/lnd/lntypes"
)

// Issues challenges in the form of invoices.
type Challenger interface {
	Challenge(price uint64) (ChallengeResult, error) // Create a challenge.
}

// A simple Challenger.
type ChallengeFactory struct {
	LightningNode
}

type ChallengeResult struct {
	lntypes.Preimage
	PaymentRequest
}

func (challenger *ChallengeFactory) Challenge(price uint64) (ChallengeResult, error) {
	secret := secrets.NewSecret()
	preimage, _ := lntypes.MakePreimage(secret[:])
	paymentRequest := InvoiceBuilder{
		RPreimage: preimage[:],
		Value:     int64(price),
		Expiry:    int64(time.Hour), // The time could change in the future
		IsKeysend: false,
		Memo:      "L402", // Idealy we would have the service name
		Private:   false,  // Not sure yet
	}
	invoice, err := challenger.LightningNode.CreateInvoice(context.Background(), paymentRequest)
	return ChallengeResult{Preimage: preimage, PaymentRequest: invoice}, err
}
