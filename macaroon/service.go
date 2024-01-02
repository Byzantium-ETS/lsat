package macaroon

import (
	"fmt"
	"strings"
)

type Tier = int8

const (
	BaseTier = 0
)

// Service represents a service associated with a macaroon caveat.
type Service struct {
	Name  string // The name of the service.
	Tier  Tier   // The tier or level of the service.
	Price uint64 // The price associated with the service.
}

func (service Service) Key() string {
	return "service"
}

func (service Service) Value() string {
	return fmt.Sprintf("%s:%d", service.Name, service.Tier)
}

func (service Service) Caveat() Caveat {
	return NewCaveat(service.Key(), service.Value())
}

func NewService(Name string, Price uint64) Service {
	return Service{Name: Name, Price: Price, Tier: BaseTier}
}

// ServiceIterator is a type representing an iterator for extracting service names
// from a sequence of caveats.
type ServiceIterator struct {
	caveats []Caveat
}

// HasNext returns true if there are more caveats and the next one is related to a service.
func (iter ServiceIterator) HasNext() bool {
	return len(iter.caveats) > 0 && iter.caveats[0].Key == "service"
}

// Next returns the next service name in the sequence of caveats.
// It assumes that there is a next service (check with HasNext before calling Next).
func (iter ServiceIterator) Next() string {
	// Split the caveat's value to extract service information.
	service := strings.Split(iter.caveats[0].Value, ":")
	iter.caveats = iter.caveats[1:]
	return service[0]
}
