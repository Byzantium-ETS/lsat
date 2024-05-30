package mock

import (
	"context"
	"encoding/hex"
	"lsat/challenge"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lntypes"
)

type TestLightningNode struct{}

func NewChallenger() challenge.Challenger {
	return &challenge.ChallengeFactory{LightningNode: &TestLightningNode{}}
}

func (ln *TestLightningNode) CreateInvoice(context.Context, challenge.CreateInvoiceRequest) (challenge.InvoiceResponse, error) {
	secret := secrets.NewSecret()

	// xor the preimage to build the invoice
	invoice := xor(secret[:])

	preimage, err := lntypes.MakePreimage(secret[:])
	if err != nil {
		return challenge.InvoiceResponse{}, err
	}

	return challenge.InvoiceResponse{
		PaymentHash: preimage.Hash(),
		Invoice:     hex.EncodeToString(invoice),
	}, nil
}

func (ln *TestLightningNode) PayInvoice(_ context.Context, invoice challenge.PayInvoiceRequest) (challenge.PayInvoiceResponse, error) {
	raw_invoice, err := hex.DecodeString(invoice.Invoice)
	if err != nil {
		return challenge.PayInvoiceResponse{}, err
	}

	preimage, err := lntypes.MakePreimage(xor(raw_invoice))
	if err != nil {
		return challenge.PayInvoiceResponse{}, err
	}

	return challenge.PayInvoiceResponse{
		PaymentId:   preimage.Hash().String(),
		Preimage:    preimage,
		PaymentHash: preimage.Hash(),
	}, nil
}

func xor(data []byte) []byte {
	maxByte := byte(255)
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ maxByte
	}
	return result
}
