package auth

import (
	"errors"
	"fmt"
	"lsat/macaroon"
	"time"
)

const (
	timeErr       = "the macaroon is expired"
	capabilityErr = "this capability is not present for this service"
)

// An interface defining methods for managing services and their capabilities.
type ServiceManager interface {
	// Services retrieves information about services with the provided names.
	Service(string) (macaroon.Service, error)

	// VerifyCaveats checks the validity of the provided caveats.
	VerifyCaveats(caveats ...macaroon.Caveat) error
}

// The configuration of every services.
type Config struct {
	services map[string]macaroon.Service
}

// Creates a new Config the provided services.
func NewConfig(services []macaroon.Service) *Config {
	serviceMap := make(map[string]macaroon.Service)
	for _, service := range services {
		serviceMap[service.Id().String()] = service
	}
	return &Config{services: serviceMap}
}

// Service retrieves information about a service with the provided name.
func (c *Config) Service(name string) (macaroon.Service, error) {
	service, exists := c.services[name]
	if !exists {
		return macaroon.Service{}, fmt.Errorf("service not found: %s", name)
	}
	return service, nil
}

// VerifyCaveats checks the validity of the provided caveats.
func (c *Config) VerifyCaveats(caveats ...macaroon.Caveat) error {
	err := c.checkExpiry(caveats...)

	if err != nil {
		return err
	}

	err = c.checkCapabilities(caveats...)

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) checkExpiry(caveats ...macaroon.Caveat) error {
	now := time.Now()

	for _, expiryTime := range macaroon.GetValue("expiry_date", caveats) {
		// Parse the value of the time caveat as a time.Time.
		expiry, err := time.Parse(time.Layout, expiryTime)

		// If there is an error parsing the time, return the error.
		if err != nil {
			return err
		}

		// Check if the expiry time is before the current time.
		if now.After(expiry) {
			return errors.New(timeErr)
		}

		now = expiry
	}

	return nil
}

func (c *Config) checkCapabilities(caveats ...macaroon.Caveat) error {
	service_id := macaroon.GetValue("service", caveats)[0]
	service := c.services[service_id]

	for _, aCapacility := range macaroon.GetValue("capability", caveats) {
		match := false
		for _, tCapability := range service.Capabilities {
			if aCapacility == tCapability {
				match = true
				break
			}
		}
		if !match {
			return errors.New(capabilityErr)
		}
	}

	return nil
}
