package tests

import (
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"testing"

	"github.com/lightningnetwork/lnd/lntypes"
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

	if err != nil {
		t.Error(err)
	}

	t.Log(preToken.Macaroon.ToJSON())

	preimage, err := lntypes.MakePreimageFromStr(preToken.InvoiceResponse.Invoice)

	if err != nil {
		t.Error(err)
	}

	lsat := macaroon.Token{
		Macaroon: preToken.Macaroon,
		Preimage: preimage,
	}

	err = minter.AuthToken(&lsat)

	if err != nil {
		t.Error(err)
	}
}
