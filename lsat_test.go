package main

import (
	"lsat/auth"
	"lsat/lightning"
	"lsat/mock"
	"testing"
)

var secretStore mock.TestStore = mock.NewTestStore()
var node mock.TestNode = mock.TestNode{}
var serviceManager mock.TestServiceManager = mock.TestServiceManager{}
var challenger lightning.ChallengeFactory = lightning.NewChallenger(&node)

func TestLsat(t *testing.T) {
	uid := secretStore.CreateUser()

	minter := auth.NewMinter(&serviceManager, &secretStore, &challenger)

	t.Log(uid)

	preToken, err := minter.MintToken(uid, mock.DogService)

	if err != nil {
		t.Error(err)
	}

	token, err := preToken.Pay(&node)

	if err != nil {
		t.Error(err)
	}

	t.Log(token.Mac)

	if err != nil {
		t.Error(err)
	}
}
