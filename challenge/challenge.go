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

// ChallengeResult represents the result of a payment challenge, including the
// preimage and payment request generated during the challenge process.
type ChallengeResult struct {
	// Preimage is the randomly generated secret that serves as the proof of payment.
	lntypes.Preimage

	// PaymentRequest is the Lightning Network invoice associated with the payment
	// challenge. Clients use this payment request to fulfill the challenge by making
	// a payment to the Lightning node.
	PaymentRequest
}

// Challenge generates a payment challenge for the specified price by creating
// a Lightning invoice and returning the associated preimage and payment request.
// The generated preimage is used to verify successful payment by the client.
//
// Parameters:
//   - price: The price of the challenge, specified in satoshis.
//
// Returns:
//   - ChallengeResult: A struct containing the preimage and payment request.
//   - error: An error, if any, encountered during the invoice creation process.
func (challenger *ChallengeFactory) Challenge(price uint64) (ChallengeResult, error) {
	// Generate a new secret and create a preimage from it.
	secret := secrets.NewSecret()
	Preimage, _ := lntypes.MakePreimage(secret[:])

	// Build an invoice with the generated preimage, price, and other details.
	invoice := InvoiceBuilder{
		RPreimage: Preimage[:],
		Value:     int64(price),
		Expiry:    int64(time.Hour), // The time could change in the future.
		IsKeysend: false,
		Memo:      "L402", // Ideally, we would have the service name.
		Private:   false,  // Not sure yet.
	}

	// Create a Lightning invoice using the built parameters.
	PaymentRequest, err := challenger.LightningNode.CreateInvoice(context.Background(), invoice)

	if err != nil {
		return ChallengeResult{}, err
	}

	// Return the ChallengeResult with the generated preimage and payment request.
	return ChallengeResult{Preimage, PaymentRequest}, nil
}
