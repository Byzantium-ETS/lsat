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

	timeKey = "time"

	timeErr     = "the macaroon is expired!"
	resourceErr = "failed to find ressource for that service!"
)

type TestServiceManager struct{}

var serviceManager TestServiceManager = TestServiceManager{}

func listCaveats(service Service) []Caveat {
	arr := make([]Caveat, 1)
	switch service.Name {
	case DogService, CatService:
		arr[0] = NewCaveat(timeKey, time.Now().Add(time.Duration(time.Hour)).Format(time.Layout))
	}
	return arr
}

func (sm *TestServiceManager) Services(cx context.Context, names ...string) ([]Service, error) {
	list := make([]Service, 0, len(names))
	for _, name := range names {
		switch name {
		case CatService:
			list = append(list, NewService(CatService, 1000))
		case DogService:
			list = append(list, NewService(DogService, 2000))
		default:
			return []Service{}, errors.New("unkown service!")
		}
	}
	return list, nil
}

func (sm *TestServiceManager) Capabilities(cx context.Context, services ...Service) ([]Caveat, error) {
	arr := make([]Caveat, 0, len(services))
	for _, service := range services {
		arr = append(arr, listCaveats(service)...)
	}
	return arr, nil
}

func (sm *TestServiceManager) VerifyCaveats(caveats ...Caveat) error {
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
	return nil
}

func (sm *TestServiceManager) GetResource(cx context.Context, macaroon Macaroon) (auth.Resource, error) {
	switch macaroon.Services().Next() {
	case DogService, CatService:
		return auth.Resource{
			Type:    "image/png",
			Content: []byte{},
		}, errors.New(resourceErr)
	}
	return auth.Resource{}, nil
}
