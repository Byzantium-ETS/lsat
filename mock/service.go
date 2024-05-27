package mock

import (
	"context"
	"errors"
	"lsat/auth"
	. "lsat/macaroon"
	"time"
)

const (
	DogService = "dogs"
	CatService = "cats"

	timeKey      = "time"
	signatureKey = "signature"

	timeErr = "the macaroon is expired"
	signErr = "the signature of that token is invalid"
)

type testServiceLimiter struct{}

func NewServiceLimiter() auth.ServiceLimiter {
	return &testServiceLimiter{}
}

// listCaveats generates a list of caveats based on the provided service.
func listCaveats(service Service) []Caveat {
	// Use a switch statement to handle different services and their respective caveats.
	switch service.Name {
	case DogService, CatService:
		// If the service is a DogService or CatService, include a time-based caveat.
		// The caveat represents an expiration time one hour from now.
		return []Caveat{NewCaveat(timeKey, time.Now().Add(time.Hour).Format(time.Layout))}
	default:
		// If the service is not explicitly handled, return an empty slice.
		return nil
	}
}

// verifyCaveats checks the validity of the provided caveats.
func (s *testServiceLimiter) VerifyCaveats(caveats ...Caveat) error {
	for _, caveat := range caveats {
		switch caveat.Key {
		case timeKey:
			// Parse the value of the time caveat as a time.Time.
			expiry, err := time.Parse(time.Layout, caveat.Value)

			// If there is an error parsing the time, return the error.
			if err != nil {
				return err
			}

			// Check if the expiry time is before the current time.
			if expiry.Before(time.Now()) {
				return errors.New(timeErr)
			}
		}
	}
	// If all checks pass, return nil (no error).
	return nil
}

func (s *testServiceLimiter) Services(cx context.Context, names ...string) ([]Service, error) {
	list := make([]Service, 0, len(names))
	for _, name := range names {
		switch name {
		case CatService:
			list = append(list, NewService(CatService, 0))
		case DogService:
			list = append(list, NewService(DogService, 0))
		default:
			return []Service{}, errors.New("unkown service")
		}
	}
	return list, nil
}

func (s *testServiceLimiter) Capabilities(cx context.Context, services ...Service) ([]Caveat, error) {
	arr := make([]Caveat, 0, len(services))
	for _, service := range services {
		arr = append(arr, listCaveats(service)...)
	}
	return arr, nil
}
