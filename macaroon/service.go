package macaroon

import (
	"fmt"
)

type Tier = int8

const (
	BaseTier = 0
)

// Un service est un caveat particulier
// C'est le premier caveat appliqué
type Service struct {
	Name  string
	Tier  Tier
	Price uint64
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
