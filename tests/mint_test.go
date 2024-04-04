package tests

import (
	"lsat/auth"
	"lsat/mock"
	"testing"
)

// I should use LndClient for testing here.
// var challenger mock.TestChallenger = mock.TestChallenger{}

func TestMintAuthMacaroon(t *testing.T) {
	serviceLimiter := mock.NewServiceLimiter()

	uid := secretStore.CreateUser()

	minter := auth.NewMinter(serviceLimiter, &secretStore, mock.NewChallenger())

	t.Log("user_id: ", uid)

	preToken, err := minter.MintToken(uid, mock.DogService)

	mac := preToken.Macaroon.String()

	t.Log("macaroon: ", mac)

	if err != nil {
		t.Error(err)
	}

	err = minter.AuthMacaroon(&preToken.Macaroon)

	if err != nil {
		t.Error(err)
	}
}
