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

var testService service.Service = service.NewService(serviceName, servicePrice)

func TestExpiryValid(t *testing.T) {
	condition := service.Expire{}

	err := condition.Satisfy(
		macaroon.NewCaveat(macaroon.ExpiryDateKey, time.Now().Add(time.Hour).Format(time.RFC3339)),
	)

	assert.Nil(t, err, err)
}

func TestExpiryInvalid(t *testing.T) {
	condition := service.Expire{}

	err := condition.Satisfy(
		macaroon.NewCaveat(macaroon.ExpiryDateKey, time.Now().Add(-time.Hour).Format(time.RFC3339)),
	)

	assert.NotNil(t, err, "Expiry should not be valid")
}

func TestCapabilitiesInvalid(t *testing.T) {
	condition := service.Capabilities{Key: "resolution"}

	err := condition.Satisfy(
		macaroon.NewCaveat("resolution", "1024x768, 800x600"),
		macaroon.NewCaveat("format", "jpeg, png"),
		macaroon.NewCaveat("resolution", "1024x768, 800x600, 1920x1080"),
	)

	assert.NotNil(t, err, "Capabilities should not be valid")
}

func TestCapabilitiesValid(t *testing.T) {
	condition := service.Capabilities{Key: "url"}

	err := condition.Satisfy(
		macaroon.NewCaveat("url", "google.com"),
		macaroon.NewCaveat("url", "facebook.com"),
	)

	assert.NotNil(t, err, "The verification should fail")
}

func TestUniqueIDValid(t *testing.T) {
	condition := service.UniqueKey{Key: "test_id"}

	err := condition.Satisfy(
		macaroon.NewCaveat("test_id", "1234"),
	)

	assert.Nil(t, err, err)
}

func TestUniqueIDInvalid(t *testing.T) {
	condition := service.UniqueKey{Key: "test_id"}

	err := condition.Satisfy(
		macaroon.NewCaveat("test_id", "1234"),
		macaroon.NewCaveat("test_id", "1234"),
	)

	assert.NotNil(t, err, "The verification should fail")
}

func TestNotBeforeValid(t *testing.T) {
	condition := service.NotBefore{}

	err := condition.Satisfy(
		macaroon.NewCaveat(macaroon.NotBeforeKey, time.Now().Add(-time.Hour).Format(time.RFC3339)),
	)

	if err != nil {
		t.Error(err)
	}
}

func TestNotBeforeInvalid1(t *testing.T) {
	condition := service.NotBefore{}

	err := condition.Satisfy(
		macaroon.NewCaveat(macaroon.NotBeforeKey, time.Now().Add(time.Hour).Format(time.RFC3339)),
	)

	assert.NotNil(t, err, "The verification should fail")
}

func TestNotBeforeInvalid2(t *testing.T) {
	condition := service.NotBefore{}

	err := condition.Satisfy(
		macaroon.NewCaveat(macaroon.NotBeforeKey, time.Now().Format(time.RFC3339)),
		macaroon.NewCaveat(macaroon.NotBeforeKey, time.Now().Add(-time.Hour).Format(time.RFC3339)),
	)

	assert.NotNil(t, err, "The verification should fail")
}

func TestGenerateIDCaveat(t *testing.T) {
	key := "test_unique_id"
	uniqueID := service.GenerateID{Name: key}

	assert.Equal(t, uniqueID.GetKey(), key, "The key should be equal to the name")
	assert.NotEmpty(t, uniqueID.GetValue(), "The value should not be empty")
	t.Log(uniqueID.GetValue())
}

func TestExpireCaveat(t *testing.T) {
	expire := service.Expire{Delay: time.Hour}
	assert.Equal(t, expire.GetKey(), macaroon.ExpiryDateKey, "The key should be equal to the expiry date key")

	future_time := expire.GetValue()
	reference_time := time.Now().Add(time.Hour).Format(time.RFC3339)

	t.Log(future_time)
	t.Log(reference_time)

	assert.Equal(t, future_time, reference_time, "The expiry date should be one hour from now")
}

func TestNotBeforeCaveat(t *testing.T) {
	notBefore := service.NotBefore{Delay: 0}
	assert.Equal(t, notBefore.GetKey(), macaroon.NotBeforeKey, "The key should be equal to the not before key")

	start_time, err := time.Parse(time.RFC3339, notBefore.GetValue())
	assert.Nil(t, err, err)

	if start_time.After(time.Now()) {
		t.Error("The start date should be in the past")
	}
}

func TestGetService(t *testing.T) {
	targetService := service.Service{
		Name:  serviceName,
		Price: servicePrice,
		Tier:  service.BaseTier,
	}

	service_id := targetService.Id()

	serviceLimiter := service.NewConfig(targetService)

	service, err := serviceLimiter.GetService(service_id)

	assert.Nil(t, err, err)
	assert.Equal(t, service, targetService)
}
