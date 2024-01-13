package tests

import (
	"lsat/auth"
	"lsat/mock"
	"testing"
)

func TestServiceAuthMacaroon(t *testing.T) {
	serviceLimiter := mock.NewServiceLimiter()

	uid := secretStore.CreateUser()

	minter := auth.NewMinter(&serviceLimiter, &secretStore, mock.NewChallenger())

	preToken, _ := minter.MintToken(uid, mock.DogService)

	mac := preToken.Macaroon

	t.Log(mac.Caveats())

	mac, err := serviceLimiter.Sign(mac)

	t.Log(mac.Caveats())

	if err != nil {
		t.Error(err)
	}

	err = serviceLimiter.VerifyMacaroon(&mac)

	if err != nil {
		t.Error(err)
	}
}
