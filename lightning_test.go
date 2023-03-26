package main

import (
	"testing"
)

const DefaultPrice uint64 = 1000

func TestPreimage(t *testing.T) {
	preimageA, _, _ := challenger.Challenge(DefaultPrice)
	preimageB, _, _ := challenger.Challenge(DefaultPrice)

	t.Log(preimageA)
	t.Log(preimageB)

	if preimageA == preimageB {
		t.Error("Two different challenges cannot have the same preimage!")
	}
}

func TestChallenge(t *testing.T) {
	preimage, _, err := challenger.Challenge(0)

	t.Log(preimage)

	if err != nil {
		t.Error(err)
	}
}

func TestInvoice(t *testing.T) {
	t.Error("Unimplemented!")
}

func TestConnection(t *testing.T) {
	t.Error("Unimplemented!")
}
