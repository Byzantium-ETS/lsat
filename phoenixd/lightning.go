package phoenixd

import (
	"context"
	"lsat/challenge"

	"github.com/lightningnetwork/lnd/lntypes"
)

type PhoenixNode struct {
	*PhoenixClient
}

func (c *PhoenixNode) CreateInvoice(_ context.Context, req challenge.CreateInvoiceRequest) (challenge.InvoiceResponse, error) {
	response, err := c.PhoenixClient.CreateInvoice(&CreateInvoiceRequest{
		Description:     req.Description,
		DescriptionHash: req.DescriptionHash.String(),
		AmountSat:       req.Amount,
		ExternalId:      req.Udata,
	})

	if err != nil {
		return challenge.InvoiceResponse{}, err
	}

	paymentHash, _ := lntypes.MakeHashFromStr(response.PaymentHash)

	return challenge.InvoiceResponse{
		PaymentHash: paymentHash,
		Invoice:     response.Serialized,
	}, nil
}

func (c *PhoenixNode) PayInvoice(_ context.Context, req challenge.PayInvoiceRequest) (challenge.PayInvoiceResponse, error) {
	response, err := c.PhoenixClient.PayInvoice(&PayInvoiceRequest{
		AmountSat: req.Amount,
		Invoice:   req.Invoice,
	})

	if err != nil {
		return challenge.PayInvoiceResponse{}, err
	}

	paymentHash, _ := lntypes.MakeHashFromStr(response.PaymentHash)
	preimage, _ := lntypes.MakePreimageFromStr(response.PaymentPreimage)

	return challenge.PayInvoiceResponse{
		PaymentId:   response.PaymentId,
		Preimage:    preimage,
		PaymentHash: paymentHash,
	}, nil
}
