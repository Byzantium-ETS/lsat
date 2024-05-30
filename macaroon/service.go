package macaroon

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Tier = int8

const (
	BaseTier = 0
)

// Service represents the configuration of a service.
type Service struct {
	Name         string        `json:"name"`         // The name of the service.
	Tier         Tier          `json:"tier"`         // The tier or level of the service.
	Price        uint64        `json:"price"`        // The price in milli-satoshi.
	Duration     time.Duration `json:"duration"`     // The lifetime of the service.
	Capabilities []Caveat      `json:"capabilities"` // The capabilities of the service.
}

// Service represents the identifiers of a Service
type ServiceId struct {
	Name string // The name of the service.
	Tier Tier   // The tier or level of the service.
}

func NewService(Name string, Price uint64) Service {
	return Service{Name: Name, Price: Price, Tier: BaseTier, Duration: time.Hour}
}

func NewServiceId(Name string, Tier Tier) ServiceId {
	return ServiceId{Name, Tier}
}

func (service ServiceId) String() string {
	return fmt.Sprintf("%s:%d", service.Name, service.Tier)
}

// Returns the identifiers of a service.
func (service *Service) Id() ServiceId {
	return ServiceId{Name: service.Name, Tier: service.Tier}
}

// The base caveats of a service.
func (service *Service) Caveats() []Caveat {
	expiry := time.Now().Round(time.Second).Add(service.Duration)
	caveats := []Caveat{
		NewCaveat("service", service.Id().String()),
		NewCaveat("expiry_date", expiry.Format(time.RFC3339)),
	}
	caveats = append(caveats, service.Capabilities...)
	return caveats
}

// ServiceIterator is a type representing an iterator for extracting service names
// from a sequence of caveats.
type ServiceIterator struct {
	caveats []Caveat
}

// HasNext returns true if there are more caveats and the next one is related to a service.
func (iter *ServiceIterator) HasNext() bool {
	return len(iter.caveats) > 0 && iter.caveats[0].Key == "service"
}

// Next returns the next service in the sequence of caveats.
func (iter *ServiceIterator) Next() ServiceId {
	// Split the caveat's value to extract service information.
	service := strings.Split(iter.caveats[0].Value, ":")
	iter.caveats = iter.caveats[1:]
	name := service[0]
	tier, _ := strconv.Atoi(service[1])
	return ServiceId{
		Name: name,
		Tier: int8(tier),
	}
}
