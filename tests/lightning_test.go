package tests

import (
	"context"
	"lsat/challenge"
	"lsat/mock"
	"lsat/secrets"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayInvoice(t *testing.T) {
	ln := mock.TestLightningNode{}

	invoice, err := ln.CreateInvoice(context.Background(), challenge.CreateInvoiceRequest{})
	if err != nil {
		t.Error(err)
	}

	payment, err := ln.PayInvoice(context.Background(), challenge.PayInvoiceRequest{Invoice: invoice.Invoice})
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, invoice.PaymentHash, payment.PaymentHash)
}

func TestXor(t *testing.T) {
	secret := secrets.NewSecret()

	xorA := xor(secret[:])
	xorB := xor(xorA)

	assert.Equal(t, secret[:], xorB)
	assert.NotEqual(t, xorA, xorB)
}

func xor(data []byte) []byte {
	maxByte := byte(255)
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ maxByte
	}
	return result
}
