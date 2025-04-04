package tests

import (
	"lsat/auth"
	"lsat/challenge"
	"lsat/mock"
	"lsat/service"
	"testing"
)

func TestMintAuthMacaroon(t *testing.T) {
	serviceLimiter := service.NewConfig(
		service.NewService(serviceName, servicePrice),
	)

	uid := secretStore.NewUser()

	minter := auth.NewMinter(serviceLimiter, secretStore, mock.NewChallenger())

	preToken, err := minter.MintToken(uid, service.NewId(serviceName, 0))

	if err != nil {
		t.Error(err)
	}

	t.Log(preToken.Macaroon.ToJSON())

	err = minter.AuthMacaroon(&preToken.Macaroon)

	if err != nil {
		t.Error(err)
	}
}

func TestMintAuthToken(t *testing.T) {
	serviceLimiter := service.NewConfig(
		service.NewService(serviceName, servicePrice),
	)

	lightningNode := mock.TestLightningNode{Balance: 1000}

	challenger := challenge.ChallengeFactory{LightningNode: &lightningNode}

	uid := secretStore.NewUser()

	minter := auth.NewMinter(serviceLimiter, secretStore, &challenger)

	preToken, err := minter.MintToken(uid, service.NewId(serviceName, 0))

	if err != nil {
		t.Error(err)
	}

	t.Log(preToken.Macaroon.ToJSON())

	token, err := preToken.Pay(&lightningNode)

	if err != nil {
		t.Error(err)
	}

	err = minter.AuthToken(&token)

	if err != nil {
		t.Error(err)
	}
}
