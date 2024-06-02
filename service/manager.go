package service

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
	Service(string) (Service, error)

	// VerifyCaveats checks the validity of the provided caveats.
	VerifyCaveats(caveats ...macaroon.Caveat) error
}

// The configuration of every services.
type Config struct {
	services map[string]Service
}

// Creates a new Config the provided services.
func NewConfig(services []Service) *Config {
	serviceMap := make(map[string]Service)
	for _, service := range services {
		serviceMap[service.Id().String()] = service
	}
	return &Config{services: serviceMap}
}

// Service retrieves information about a service with the provided name.
func (c *Config) Service(name string) (Service, error) {
	service, exists := c.services[name]
	if !exists {
		return Service{}, fmt.Errorf("service not found: %s", name)
	}
	return service, nil
}

// VerifyCaveats checks the validity of the provided caveats.
func (c *Config) VerifyCaveats(caveats ...macaroon.Caveat) error {
	err := c.checkExpiry(caveats...)

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) checkExpiry(caveats ...macaroon.Caveat) error {
	now := time.Now()
	var previousExpiry time.Time

	for i, expiryTime := range macaroon.GetValue(macaroon.ExpiryDateKey, caveats) {
		// Parse the value of the time caveat as a time.Time.
		expiry, err := time.Parse(time.RFC3339, expiryTime)

		// If there is an error parsing the time, return the error.
		if err != nil {
			return err
		}

		if i == 0 {
			// The first expiry_date should be after now.
			if now.After(expiry) {
				return errors.New(timeErr)
			}
		} else {
			// Each following expiry_date should be more strict or before the previous expiry date.
			if expiry.After(previousExpiry) {
				return fmt.Errorf("%s is not more strict than the previous one", macaroon.ExpiryDateKey)
			}
		}

		// Update previousExpiry to the current expiry.
		previousExpiry = expiry
	}

	// now must be before all the expiry_date.
	if now.After(previousExpiry) {
		return errors.New(timeErr)
	}

	return nil
}
