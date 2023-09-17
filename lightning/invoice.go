package lightning

import (
	"context"
	"lsat/secrets"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
)

type InvoiceBuilder = lnrpc.Invoice     /// A type used to build an Invoice.
type Invoice = lnrpc.AddInvoiceResponse /// A BOLT11 invoice.

type InvoiceHandler interface { // Je ne suis pas encore sur où est le bon endroit pour introduire les cxs.
	Pay(context.Context, Invoice) (lntypes.Preimage, error)         // Pay a BOLT11 invoice.
	CreateInvoice(context.Context, InvoiceBuilder) (Invoice, error) // Create a BOLT11 invoice.
}

// Issues challenges in the form of invoices.
type Challenger interface {
	Challenge(price uint64) (lntypes.Preimage, lnrpc.AddInvoiceResponse, error)
}

// A simple Challenger.
type ChallengeFactory struct {
	InvoiceHandler InvoiceHandler
}

func NewChallenger(InvoiceHandler InvoiceHandler) ChallengeFactory {
	return ChallengeFactory{InvoiceHandler}
}

func (challenger *ChallengeFactory) Challenge(price uint64) (lntypes.Preimage, Invoice, error) {
	secret := secrets.NewSecret()
	preimage, _ := lntypes.MakePreimage(secret[:])
	pr := InvoiceBuilder{
		RPreimage: preimage[:],
		Value:     int64(price),
		Expiry:    int64(time.Hour), // The time could change in the future
		IsKeysend: false,
		Memo:      "L402", // Idealy we would have the service name
		Private:   false,  // Not sure yet
	}
	invoice, err := challenger.InvoiceHandler.CreateInvoice(context.Background(), pr)
	return preimage, invoice, err
}

func (challenger *ChallengeFactory) Pay(invoice Invoice) (lntypes.Preimage, error) {
	return challenger.InvoiceHandler.Pay(context.Background(), invoice)
}

func (challenger *ChallengeFactory) CreateInvoice(pr InvoiceBuilder) (Invoice, error) {
	return challenger.InvoiceHandler.CreateInvoice(context.Background(), pr)
}
