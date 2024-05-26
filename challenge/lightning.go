package challenge

import (
	"context"

	"github.com/lightningnetwork/lnd/lntypes"
)

type CreateInvoiceRequest struct {
	Description     string
	DescriptionHash lntypes.Hash
	Amount          uint64
	Udata           any
}

type PayInvoiceRequest struct {
	Amount  uint64
	Invoice string
}

type InvoiceResponse struct {
	Preimage    lntypes.Preimage
	PaymentHash lntypes.Hash
	Invoice     string
}

type PayInvoiceResponse struct {
	PaymentId   string
	Preimage    lntypes.Preimage
	PaymentHash lntypes.Hash
}

// A Lightning Network node.
type LightningNode interface {
	// PayInvoice sends a payment using the Lightning Network.
	PayInvoice(context.Context, PayInvoiceRequest) (PayInvoiceResponse, error)

	// CreateInvoice creates an invoice for receiving payments on the Lightning Network.
	CreateInvoice(context.Context, CreateInvoiceRequest) (InvoiceResponse, error)
}
