package mock

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"lsat/challenge"
	"lsat/secrets"
	"math"

	"github.com/lightningnetwork/lnd/lntypes"
)

type TestLightningNode struct {
	Balance uint64
}

type invoice struct {
	RawInvoice []byte `json:"raw_invoice"`
	Amount     uint64 `json:"amount"`
}

func NewChallenger() challenge.Challenger {
	return &challenge.ChallengeFactory{LightningNode: &TestLightningNode{Balance: math.MaxUint64}}
}

func (ln *TestLightningNode) CreateInvoice(ctx context.Context, req challenge.CreateInvoiceRequest) (challenge.InvoiceResponse, error) {
	secret := secrets.NewSecret()

	// xor the preimage to build the invoice
	rawInvoice := xor(secret[:])

	inv := invoice{
		RawInvoice: rawInvoice,
		Amount:     req.Amount,
	}

	invoiceJSON, err := json.Marshal(inv)
	if err != nil {
		return challenge.InvoiceResponse{}, err
	}

	preimage, err := lntypes.MakePreimage(secret[:])
	if err != nil {
		return challenge.InvoiceResponse{}, err
	}

	return challenge.InvoiceResponse{
		PaymentHash: preimage.Hash(),
		Invoice:     hex.EncodeToString(invoiceJSON),
	}, nil
}

func (ln *TestLightningNode) PayInvoice(ctx context.Context, req challenge.PayInvoiceRequest) (challenge.PayInvoiceResponse, error) {
	invoiceJSON, err := hex.DecodeString(req.Invoice)
	if err != nil {
		return challenge.PayInvoiceResponse{}, err
	}

	var inv invoice
	err = json.Unmarshal(invoiceJSON, &inv)
	if err != nil {
		return challenge.PayInvoiceResponse{}, err
	}

	if inv.Amount > ln.Balance {
		return challenge.PayInvoiceResponse{}, errors.New("insufficient balance")
	}

	preimage, err := lntypes.MakePreimage(xor(inv.RawInvoice))
	if err != nil {
		return challenge.PayInvoiceResponse{}, err
	}

	ln.Balance -= inv.Amount

	return challenge.PayInvoiceResponse{
		PaymentId:   preimage.Hash().String(),
		Preimage:    preimage,
		PaymentHash: preimage.Hash(),
	}, nil
}

func xor(data []byte) []byte {
	maxByte := byte(math.MaxUint8)
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ maxByte
	}
	return result
}
