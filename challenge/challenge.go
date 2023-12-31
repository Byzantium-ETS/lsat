package challenge

import (
	"context"
	"lsat/secrets"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
)

type InvoiceBuilder = lnrpc.Invoice            /// A type used to build an Invoice.
type PaymentRequest = lnrpc.AddInvoiceResponse /// A BOLT11 invoice.

type InvoiceHandler interface { // Je ne suis pas encore sur où est le bon endroit pour introduire les cxs.
	SendPayment(context.Context, PaymentRequest) (lntypes.Preimage, error) // Pay a BOLT11 invoice. The server should not be required to pay.
	CreateInvoice(context.Context, InvoiceBuilder) (PaymentRequest, error) // Create a BOLT11 invoice.
}

// Issues challenges in the form of invoices.
type Challenger interface {
	Challenge(price uint64) (ChallengeResult, error)
}

// A simple Challenger.
type ChallengeFactory struct {
	InvoiceHandler InvoiceHandler
}

func NewChallenger(InvoiceHandler InvoiceHandler) ChallengeFactory {
	return ChallengeFactory{InvoiceHandler}
}

type ChallengeResult struct {
	Preimage lntypes.Preimage
	Invoice  PaymentRequest
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
	invoice, err := challenger.InvoiceHandler.CreateInvoice(context.Background(), paymentRequest)
	return ChallengeResult{Preimage: preimage, Invoice: invoice}, err
}
