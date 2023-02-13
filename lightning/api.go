package lightning

import "time"

type Multiplier int

const (
	Milli = 1
	Micro = 1
	Nano  = 1
	Pico  = 1
)

type SatoshiAmount struct {
	raw_amount int
	unit       Multiplier
}

func MilliSats(amount int) SatoshiAmount {
	return SatoshiAmount{raw_amount: amount, unit: Milli}
}

func (sats SatoshiAmount) Unwrap() int {
	return sats.raw_amount * int(sats.unit)
}

type PaymentRequest struct {
	add_index    uint64
	invoice      string
	r_hash       []uint8
	payment_addr []uint8
}

type LightningNode interface {
	CreateInvoice(value_msat SatoshiAmount, expiry time.Time, private bool, memo string, preimage string) (PaymentRequest, error)
	SubscribeInvoice(r_hash []uint8) error
}
