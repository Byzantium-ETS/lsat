package tests

import (
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"reflect"
	"testing"
)

var service macaroon.Service = macaroon.NewService("image", 1000)

var caveat macaroon.Caveat = macaroon.NewCaveat("expiry", "12:00 PM")

func TestServiceAuthMacaroon(t *testing.T) {
	serviceLimiter := mock.NewServiceLimiter()

	uid := secretStore.CreateUser()

	minter := auth.NewMinter(&serviceLimiter, &secretStore, mock.NewChallenger())

	preToken, _ := minter.MintToken(uid, mock.DogService)

	mac := preToken.Macaroon

	t.Log(mac)

	mac, err := serviceLimiter.Sign(mac)

	t.Log(mac)

	if err != nil {
		t.Error(err)
	}

	err = serviceLimiter.VerifyMacaroon(&mac)

	if err != nil {
		t.Error(err)
	}
}

func TestServiceAuthMacaroonSignature(t *testing.T) {
	serviceLimiter := mock.NewServiceLimiter()

	uid := secretStore.CreateUser()

	minter := auth.NewMinter(&serviceLimiter, &secretStore, mock.NewChallenger())

	preToken, _ := minter.MintToken(uid, mock.DogService)

	mac := preToken.Macaroon

	signedMac, err := serviceLimiter.Sign(mac)

	if err != nil {
		t.Error(err)
	}

	newMac, err := macaroon.DecodeBase64(signedMac.String())

	if err != nil {
		t.Error(err)
	}

	t.Log(signedMac)
	t.Log(newMac)

	if !reflect.DeepEqual(signedMac, newMac) {
		t.Error("the macaroons should be identical!")
	}
}

func TestServiceAuthMacaroonEncoded(t *testing.T) {
	serviceLimiter := mock.NewServiceLimiter()

	uid := secretStore.CreateUser()

	minter := auth.NewMinter(&serviceLimiter, &secretStore, mock.NewChallenger())

	preToken, _ := minter.MintToken(uid, mock.DogService)

	mac := preToken.Macaroon

	mac, err := serviceLimiter.Sign(mac)

	// t.Log(mac)

	if err != nil {
		t.Error(err)
	}

	mac, err = macaroon.DecodeBase64(mac.String())

	if err != nil {
		t.Error(err)
	}

	err = serviceLimiter.VerifyMacaroon(&mac)

	if err != nil {
		t.Error(err)
	}
}
