package service

import (
	"fmt"
	"lsat/macaroon"
	"strconv"
	"strings"
	"time"
)

type Tier = int8

const (
	BaseTier Tier = 0
)

// A callback function for the service.
type TokenCallback func(any) error

// Service represents the configuration of a service.
type Service struct {
	Name              string            // The name of the service.
	Tier              Tier              // The tier or level of the service.
	Price             uint64            // The price in milli-satoshi.
	Duration          time.Duration     // The lifetime of the service.
	FirstPartyCaveats []macaroon.Caveat // The caveats of the service.
	Conditions        []Condition       // The conditions of the service.
	Callback          TokenCallback     // The callback function for the service.
}

// Service represents the identifiers of a Service
type ServiceID struct {
	Name string // The name of the service.
	Tier Tier   // The tier or level of the service.
}

// Create a new service configuration.
func NewService(Name string, Price uint64) Service {
	return Service{
		Name:       Name,
		Price:      Price,
		Tier:       BaseTier,
		Duration:   time.Hour,
		Conditions: []Condition{Timeout{}},
	}
}

func ParseServiceID(serviceStr string) (ServiceID, error) {
	parts := strings.Split(serviceStr, ":")
	if len(parts) != 2 {
		return ServiceID{}, fmt.Errorf("invalid service ID format: %s", serviceStr)
	}

	tier, err := strconv.Atoi(parts[1])
	if err != nil {
		return ServiceID{}, fmt.Errorf("invalid tier: %s", parts[1])
	}

	return ServiceID{Name: parts[0], Tier: Tier(tier)}, nil
}

// Create a new service identifier.
func NewId(Name string, Tier Tier) ServiceID {
	return ServiceID{Name, Tier}
}

func (service ServiceID) String() string {
	return fmt.Sprintf("%s:%d", service.Name, service.Tier)
}

// Returns the identifiers of a service.
func (service *Service) Id() ServiceID {
	return ServiceID{Name: service.Name, Tier: service.Tier}
}

// The base caveats of a service.
func (service *Service) Caveats() []macaroon.Caveat {
	expiry := time.Now().Add(service.Duration)
	caveats := []macaroon.Caveat{
		macaroon.NewCaveat(macaroon.ServiceKey, service.Id().String()),
		macaroon.NewCaveat(macaroon.ExpiryDateKey, expiry.Format(time.RFC3339)),
	}
	caveats = append(caveats, service.FirstPartyCaveats...)
	return caveats
}
