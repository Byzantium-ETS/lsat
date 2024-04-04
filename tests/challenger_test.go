package tests

import (
	"lsat/mock"
	"testing"
)

const defaultPrice uint64 = 1000

func TestPreimage(t *testing.T) {
	// Should be replaced by LndClient
	var challenger = mock.NewChallenger()

	resultA, _ := challenger.Challenge(defaultPrice)
	resultB, _ := challenger.Challenge(defaultPrice)

	t.Log(resultA.Preimage)
	t.Log(resultB.Preimage)

	if resultA.Preimage == resultB.Preimage {
		t.Error("Two different challenges cannot have the same preimage!")
	}
}

func TestChallenge(t *testing.T) {
	// Should be replaced by LndClient
	var challenger = mock.NewChallenger()

	result, err := challenger.Challenge(0)

	t.Log(result.Preimage)

	if err != nil {
		t.Error(err)
	}
}

func TestPaymentRequest(t *testing.T) {
	// Should be replaced by LndClient
	var challenger = mock.NewChallenger()

	result, _ := challenger.Challenge(0)

	t.Log(result.GetPaymentRequest())

	if result.Preimage.String() != result.GetPaymentRequest() {
		t.Error("the test challenger should produce an invoice identical to the preimage!")
	}
}

func TestInvoice(t *testing.T) {
	t.Error("Unimplemented!")
}

func TestConnection(t *testing.T) {
	t.Error("Unimplemented!")
}
