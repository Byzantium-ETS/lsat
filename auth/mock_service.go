package auth

import (
	"context"
	"errors"
	"time"
)

const (
	dogService = "dogs"
	catService = "cats"

	timeKey = "time"

	timeErr = "the macaroon is expired!"
)

type TestServiceManager struct {
}

func listCaveats(service Service) []Caveat {
	arr := make([]Caveat, 1)
	switch service.Name {
	case dogService, catService:
		arr = append(arr, NewCaveat("time", time.Now().Add(time.Duration(1000)).Format(time.Layout)))
	}
	return arr
}

func (sm *TestServiceManager) Services(cx context.Context, names ...string) ([]Service, error) {
	list := make([]Service, len(names))
	for _, name := range names {
		switch name {
		case "cat":
			list = append(list, NewService("cats", 1000))
		case "dog":
			list = append(list, NewService("dogs", 2000))
		default:
			return []Service{}, errors.New("unkown service!")
		}
	}
	return list, nil
}

func (sm *TestServiceManager) Capabilities(cx context.Context, services ...Service) ([]Caveat, error) {
	arr := make([]Caveat, 10)
	for _, service := range services {
		arr = append(arr, listCaveats(service)...)
	}
	return arr, nil
}

func (sm *TestServiceManager) VerifyCaveats(service Service, caveats ...Caveat) error {
	for _, caveat := range caveats {
		switch caveat.Key {
		case timeKey:
			expiry, err := time.Parse(time.Layout, caveat.Value)

			if err != nil {
				return err
			}

			if expiry.Before(time.Now()) {
				return errors.New(timeErr)
			}
		}
	}
}
