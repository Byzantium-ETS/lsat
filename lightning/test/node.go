package test

import (
	"lsat/lightning"
	"time"
)

type Node struct {
}

func (node Node) CreateInvoice(value_msat lightning.SatoshiAmount, expiry time.Time, private bool, memo string, preimage string) lightning.PaymentRequest {

}
