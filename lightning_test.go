package main

import (
	"lsat/mock"
	"testing"
)

const DefaultPrice int64 = 1000

func TestPreimage(t *testing.T) {
	challenger := mock.NewTestChallenger()

	preimageA, _, _ := challenger.Challenge(DefaultPrice)
	preimageB, _, _ := challenger.Challenge(DefaultPrice)

	if preimageA == preimageB {
		t.Error("Two different challenges cannot have the same preimage")
	}
}
