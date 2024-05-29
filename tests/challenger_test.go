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

	t.Log(resultA.PaymentHash)
	t.Log(resultB.PaymentHash)

	if resultA.PaymentHash == resultB.PaymentHash {
		t.Error("Two different challenges cannot have the same preimage!")
	}
}

func TestChallenge(t *testing.T) {
	// Should be replaced by LndClient
	var challenger = mock.NewChallenger()

	result, err := challenger.Challenge(0)

	t.Log(result.PaymentHash)

	if err != nil {
		t.Error(err)
	}
}
