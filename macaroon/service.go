package macaroon

import (
	"errors"
	"fmt"
	"strings"
)

type Tier = int8

const (
	ServiceCaveatKey = "services"
)

// Un service est un caveat particulier
// C'est le premier caveat appliqué
type Service struct {
	Name  string
	Tier  Tier
	Price uint64
}

func NewService(Name string, Price uint64) Service {
	return Service{Name: Name, Price: Price, Tier: 0}
}

func FmtServices(services ...Service) (string, error) {
	var s strings.Builder

	fmt.Fprintf(&s, ServiceCaveatKey)
	for _, service := range services {
		if service.Name == "" {
			return "", errors.New("missing service name!")
		}

		fmt.Fprintf(&s, service.Name+":"+string(service.Tier))
	}
	if s.Len() == 0 {
		return "", errors.New("no services!")
	} else {
		return s.String(), nil
	}
}
