package service

import (
	"fmt"
	"lsat/macaroon"
)

const (
	timeErr       = "the macaroon is expired"
	capabilityErr = "this capability is not present for this service"
)

// An interface defining methods for managing services and their capabilities.
type ServiceManager interface {
	// GetServices retrieves information about services with the provided names.
	GetService(ServiceID) (Service, error)

	// VerifyCaveats checks the validity of the provided caveats.
	VerifyCaveats(caveats ...macaroon.Caveat) error
}

// The configuration of every services.
type Config struct {
	services map[ServiceID]Service
}

// Creates a new Config the provided services.
func NewConfig(services ...Service) *Config {
	serviceMap := make(map[ServiceID]Service)
	for _, service := range services {
		serviceMap[service.Id()] = service
	}
	return &Config{services: serviceMap}
}

// GetService retrieves information about a service with the provided name.
func (c *Config) GetService(id ServiceID) (Service, error) {
	service, exists := c.services[id]
	if !exists {
		return Service{}, fmt.Errorf("service not found: %s", id.String())
	}
	return service, nil
}

// VerifyCaveats checks the validity of the provided caveats.
func (c *Config) VerifyCaveats(caveats ...macaroon.Caveat) error {
	iter := macaroon.NewIterator(macaroon.ServiceKey, caveats)
	for iter.HasNext() {
		service_str := iter.Next()
		service_id, _ := ParseServiceID(service_str)
		service, _ := c.services[service_id]
		for _, condition := range service.Conditions {
			err := condition.Satisfy(caveats...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
