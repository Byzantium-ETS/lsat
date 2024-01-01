package challenge

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
)

type InvoiceBuilder = lnrpc.Invoice            /// A type used to build an Invoice.
type PaymentRequest = lnrpc.AddInvoiceResponse /// A BOLT11 invoice.

// LightningNode defines the interface for a Lightning Network node.
type LightningNode interface {
	// SendPayment sends a payment using the Lightning Network.
	SendPayment(context.Context, PaymentRequest) (lntypes.Preimage, error)

	// CreateInvoice creates an invoice for receiving payments on the Lightning Network.
	CreateInvoice(context.Context, InvoiceBuilder) (PaymentRequest, error)
}
