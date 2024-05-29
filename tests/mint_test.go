package tests

import (
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"testing"

	"github.com/lightningnetwork/lnd/lntypes"
)

func TestMintAuthMacaroon(t *testing.T) {
	serviceLimiter := auth.NewServiceManager([]macaroon.Service{
		macaroon.NewService(serviceName, servicePrice),
	})

	uid := secretStore.NewUser()

	minter := auth.NewMinter(serviceLimiter, &secretStore, mock.NewChallenger())

	preToken, err := minter.MintToken(uid, serviceName)

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
	serviceLimiter := auth.NewServiceManager([]macaroon.Service{
		macaroon.NewService(serviceName, servicePrice),
	})

	uid := secretStore.NewUser()

	minter := auth.NewMinter(serviceLimiter, &secretStore, mock.NewChallenger())

	preToken, err := minter.MintToken(uid, serviceName)

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
