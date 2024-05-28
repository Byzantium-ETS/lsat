package tests

import (
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"testing"
)

func TestMintAuthMacaroon(t *testing.T) {
	serviceLimiter := mock.NewServiceLimiter()

	uid := secretStore.NewUser()

	minter := auth.NewMinter(serviceLimiter, &secretStore, mock.NewChallenger())

	preToken, err := minter.MintToken(uid, mock.DogService)

	t.Log(preToken.Macaroon.ToJSON())

	if err != nil {
		t.Error(err)
	}

	err = minter.AuthMacaroon(&preToken.Macaroon)

	if err != nil {
		t.Error(err)
	}
}

func TestMintAuthToken(t *testing.T) {
	serviceLimiter := mock.NewServiceLimiter()

	uid := secretStore.NewUser()

	minter := auth.NewMinter(serviceLimiter, &secretStore, mock.NewChallenger())

	preToken, err := minter.MintToken(uid, mock.DogService)

	t.Log(preToken.Macaroon.ToJSON())

	if err != nil {
		t.Error(err)
	}

	lsat := macaroon.Token{
		Macaroon: preToken.Macaroon,
		Preimage: preToken.InvoiceResponse.Preimage,
	}

	err = minter.AuthToken(&lsat)

	if err != nil {
		t.Error(err)
	}
}
