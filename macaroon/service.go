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
	Name         string        // The name of the service.
	Tier         Tier          // The tier or level of the service.
	Price        uint64        // The price in milli-satoshi.
	Duration     time.Duration // The lifetime of the service.
	Capabilities []string      // The capabilities of the service.
}

// Service represents the identifiers of a Service
type ServiceId struct {
	Name string // The name of the service.
	Tier Tier   // The tier or level of the service.
}

func NewService(Name string, Price uint64) Service {
	return Service{Name: Name, Price: Price, Tier: BaseTier}
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
	expiry := time.Now().Add(service.Duration)
	caveats := []Caveat{
		NewCaveat("service", service.Id().String()),
		NewCaveat("expiry_date", expiry.String()),
	}
	for _, capability := range service.Capabilities {
		caveats = append(caveats, NewCaveat("capability", capability))
	}
	return caveats
}

// ServiceIterator is a type representing an iterator for extracting service names
// from a sequence of caveats.
type ServiceIterator struct {
	caveats []Caveat
}

// func (iter ServiceIterator) String() string {
// 	s := "services: "

// 	for i := 0; iter.caveats[i].Key == "service"; i++ {
// 		s += iter.caveats[i].Value + ","
// 	}

// 	return s
// }

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
