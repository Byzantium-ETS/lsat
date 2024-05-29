package auth

import (
	"context"
	"errors"
	"fmt"
	"lsat/macaroon"
	"time"
)

const (
	timeErr       = "the macaroon is expired"
	capabilityErr = "this capability is not present for this service"
)

// ServiceLimiter is an interface defining methods for managing services and their capabilities.
type ServiceLimiter interface {
	// Services retrieves information about services with the provided names.
	Service(context.Context, string) (macaroon.Service, error)

	// VerifyCaveats checks the validity of the provided caveats.
	VerifyCaveats(caveats ...macaroon.Caveat) error
}

// ServiceManager implements the ServiceLimiter interface and manages multiple services.
type ServiceManager struct {
	services map[string]macaroon.Service
}

// NewServiceManager creates a new ServiceManager with the provided services.
func NewServiceManager(services []macaroon.Service) *ServiceManager {
	serviceMap := make(map[string]macaroon.Service)
	for _, service := range services {
		serviceMap[service.Name] = service
	}
	return &ServiceManager{services: serviceMap}
}

// Service retrieves information about a service with the provided name.
func (sm *ServiceManager) Service(ctx context.Context, name string) (macaroon.Service, error) {
	service, exists := sm.services[name]
	if !exists {
		return macaroon.Service{}, fmt.Errorf("service not found: %s", name)
	}
	return service, nil
}

// VerifyCaveats checks the validity of the provided caveats.
func (sm *ServiceManager) VerifyCaveats(caveats ...macaroon.Caveat) error {
	err := sm.checkExpiry(caveats...)

	if err != nil {
		return err
	}

	err = sm.checkCapabilities(caveats...)

	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) checkExpiry(caveats ...macaroon.Caveat) error {
	timeLimit := time.Now()

	for _, expiryTime := range macaroon.GetValue("expiry", caveats) {
		// Parse the value of the time caveat as a time.Time.
		expiry, err := time.Parse(time.Layout, expiryTime)

		// If there is an error parsing the time, return the error.
		if err != nil {
			return err
		}

		// Check if the expiry time is before the current time.
		if expiry.Before(timeLimit) {
			return errors.New(timeErr)
		}

		timeLimit = expiry
	}

	return nil
}

func (sm *ServiceManager) checkCapabilities(caveats ...macaroon.Caveat) error {
	service_id := macaroon.GetValue("service", caveats)[0]
	service := sm.services[service_id]

	for _, aCapacility := range macaroon.GetValue("capabilities", caveats) {
		for _, tCapability := range service.Capabilities {
			if aCapacility == tCapability {
				break
			}
		}
		return errors.New(capabilityErr)
	}

	return nil
}
