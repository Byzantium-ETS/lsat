package tests

import (
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	serviceName  = "image"
	servicePrice = 1000
)

var service macaroon.Service = macaroon.NewService("image", 1000)

var caveat macaroon.Caveat = macaroon.NewCaveat("expiry", "12:00 PM")

func TestVerifyCaveats(t *testing.T) {
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

	preToken, err := minter.MintToken(uid, serviceName+":0")

	if err != nil {
		t.Error(err)
	}

	mac := preToken.Macaroon

	t.Log(mac.ToJSON())

	err = serviceLimiter.VerifyCaveats(mac.Caveats()...)

	if err != nil {
		t.Error(err)
	}
}

func TestService(t *testing.T) {
	targetService := macaroon.Service{
		Name:     serviceName,
		Price:    servicePrice,
		Tier:     macaroon.BaseTier,
		Duration: time.Hour,
	}

	service_id := targetService.Id().String()

	serviceLimiter := auth.NewServiceManager([]macaroon.Service{targetService})

	service, err := serviceLimiter.Service(service_id)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, service, targetService)
}
