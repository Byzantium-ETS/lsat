package tests

import (
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"testing"
	"time"
)

const (
	serviceName  = "image"
	servicePrice = 1000
)

var service macaroon.Service = macaroon.NewService("image", 1000)

var caveat macaroon.Caveat = macaroon.NewCaveat("expiry", "12:00 PM")

func TestServiceAuthMacaroon(t *testing.T) {
	serviceLimiter := auth.NewServiceManager([]macaroon.Service{
		{
			Name:     serviceName,
			Price:    servicePrice,
			Tier:     macaroon.BaseTier,
			Duration: time.Hour,
		},
	})

	uid := secretStore.NewUser()

	minter := auth.NewMinter(serviceLimiter, &secretStore, mock.NewChallenger())

	preToken, _ := minter.MintToken(uid, serviceName)

	mac := preToken.Macaroon

	t.Log(mac)

	err := serviceLimiter.VerifyCaveats(mac.Caveats()...)

	if err != nil {
		t.Error(err)
	}
}

func TestServiceAuthMacaroonEncoded(t *testing.T) {
	serviceLimiter := auth.NewServiceManager([]macaroon.Service{
		{
			Name:     serviceName,
			Price:    servicePrice,
			Tier:     macaroon.BaseTier,
			Duration: time.Hour,
		},
	})

	uid := secretStore.NewUser()

	minter := auth.NewMinter(serviceLimiter, &secretStore, mock.NewChallenger())

	preToken, _ := minter.MintToken(uid, serviceName)

	mac := preToken.Macaroon

	mac, err := macaroon.DecodeBase64(mac.String())

	if err != nil {
		t.Error(err)
	}

	err = serviceLimiter.VerifyCaveats(mac.Caveats()...)

	if err != nil {
		t.Error(err)
	}
}
