package main

import (
	"lsat/auth"
	"lsat/mock"
	"testing"
)

var secretStore mock.TestStore = mock.NewTestStore()
var serviceManager mock.TestServiceManager = mock.TestServiceManager{}

// I should use LndClient for testing here.
// var challenger mock.TestChallenger = mock.TestChallenger{}

func TestAuthMacaroon(t *testing.T) {
	uid := secretStore.CreateUser()

	minter := auth.NewMinter(&serviceManager, &secretStore, &mock.TestChallenger{})

	t.Log("user_id: ", uid)

	preToken, err := minter.MintToken(uid, mock.DogService)

	mac := preToken.Macaroon.String()

	t.Log("macaroon: ", mac)

	if err != nil {
		t.Error(err)
	}

	err = minter.AuthMacaroon(&preToken.Macaroon)

	if err != nil {
		t.Error(err)
	}
}
