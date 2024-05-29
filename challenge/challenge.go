package challenge

import (
	"context"
)

// Issues challenges in the form of invoices.
type Challenger interface {
	// The price is in satoshi.
	Challenge(price uint64) (InvoiceResponse, error) // Create a challenge.
}

// A simple Challenger.
type ChallengeFactory struct {
	LightningNode
}

// Challenge generates a payment challenge for the specified price by creating a Lightning invoice
func (challenger *ChallengeFactory) Challenge(price uint64) (InvoiceResponse, error) {
	// Build an invoice with the generated preimage, price, and other details.
	invoice := CreateInvoiceRequest{
		Amount:      uint64(price),
		Description: "L402", // Ideally, we would have the service name.
	}

	// Create a Lightning invoice using the built parameters.
	response, err := challenger.LightningNode.CreateInvoice(context.Background(), invoice)

	if err != nil {
		return InvoiceResponse{}, err
	}

	// Return the ChallengeResult with the generated preimage and payment request.
	return response, nil
}
