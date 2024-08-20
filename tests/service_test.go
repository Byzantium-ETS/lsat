package tests

import (
	"lsat/macaroon"
	"lsat/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	serviceName  = "image"
	servicePrice = 1000
)

var testService service.Service = service.NewService("image", 1000)

func TestExpiryValid(t *testing.T) {
	serviceLimiter := service.NewConfig([]service.Service{
		{
			Name:     serviceName,
			Price:    servicePrice,
			Tier:     service.BaseTier,
			Duration: time.Hour,
			Conditions: []service.Condition{
				service.Timeout{},
			},
		},
	})

	caveat := macaroon.NewCaveat(macaroon.ExpiryDateKey, time.Now().Add(time.Hour).Format(time.RFC3339))

	err := serviceLimiter.VerifyCaveats(macaroon.NewCaveat("service", testService.Id().String()), caveat)

	if err != nil {
		t.Error(err)
	}
}

func TestExpiryInvalid(t *testing.T) {
	serviceLimiter := service.NewConfig([]service.Service{
		{
			Name:     serviceName,
			Price:    servicePrice,
			Tier:     service.BaseTier,
			Duration: time.Hour,
			Conditions: []service.Condition{
				service.Timeout{},
			},
		},
	})

	caveat := macaroon.NewCaveat(macaroon.ExpiryDateKey, time.Now().Add(-time.Hour).Format(time.RFC3339))

	t.Log(caveat)

	err := serviceLimiter.VerifyCaveats(macaroon.NewCaveat("service", testService.Id().String()), caveat)

	if err == nil {
		t.Error("Expiry should not be valid")
	}
}

func TestCapabilities(t *testing.T) {
	serviceLimiter := service.NewConfig([]service.Service{
		{
			Name:     serviceName,
			Price:    servicePrice,
			Tier:     service.BaseTier,
			Duration: time.Hour,
			Conditions: []service.Condition{
				service.Capabilities{Key: "resolution"},
			},
		},
	})

	err := serviceLimiter.VerifyCaveats(
		macaroon.NewCaveat("service", testService.Id().String()),
		macaroon.NewCaveat("resolution", "1024x768, 800x600"),
		macaroon.NewCaveat("format", "jpeg, png"),
		macaroon.NewCaveat("resolution", "1024x768, 800x600, 1920x1080"),
	)

	if err != nil {
		t.Log(err)
	} else {
		t.Error("Capabilities should not be valid")
	}
}

func TestUniqueKey(t *testing.T) {
	serviceLimiter := service.NewConfig([]service.Service{
		{
			Name:     serviceName,
			Price:    servicePrice,
			Tier:     service.BaseTier,
			Duration: time.Hour,
			Conditions: []service.Condition{
				service.Capabilities{Key: "url"},
			},
		},
	})

	err := serviceLimiter.VerifyCaveats(
		macaroon.NewCaveat("service", testService.Id().String()),
		macaroon.NewCaveat("url", "google.com"),
		macaroon.NewCaveat("url", "facebook.com"),
	)

	if err != nil {
		t.Log(err)
	} else {
		t.Error("The verification should fail")
	}
}

func TestService(t *testing.T) {
	targetService := service.Service{
		Name:     serviceName,
		Price:    servicePrice,
		Tier:     service.BaseTier,
		Duration: time.Hour,
	}

	service_id := targetService.Id()

	serviceLimiter := service.NewConfig([]service.Service{targetService})

	service, err := serviceLimiter.GetService(service_id)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, service, targetService)
}
