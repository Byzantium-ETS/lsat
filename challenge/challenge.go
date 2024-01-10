package challenge

import (
	"context"
	"lsat/secrets"
	"time"

	"github.com/lightningnetwork/lnd/lntypes"
)

// Issues challenges in the form of invoices.
type Challenger interface {
	// The price is in satoshi.
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
	Preimage, _ := lntypes.MakePreimage(secret[:])

	invoice := InvoiceBuilder{
		RPreimage: Preimage[:],
		Value:     int64(price),
		Expiry:    int64(time.Hour), // The time could change in the future
		IsKeysend: false,
		Memo:      "L402", // Idealy we would have the service name
		Private:   false,  // Not sure yet
	}

	PaymentRequest, err := challenger.LightningNode.CreateInvoice(context.Background(), invoice)

	if err != nil {
		return ChallengeResult{}, err
	}

	return ChallengeResult{Preimage, PaymentRequest}, nil
}
