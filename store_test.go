package main

import (
	"lsat/macaroon"
	"lsat/mock"
	"testing"
)

var service macaroon.Service = macaroon.NewService("images", 1000)

var caveats []macaroon.Caveat = []macaroon.Caveat{
	macaroon.NewCaveat("image", "test.png"),
}

func TestSecret(t *testing.T) {
	store := mock.NewTestStore()

	uid := store.CreateUser()

	asecret, err := store.Secret(uid)

	if err != nil {
		t.Error("Secret not found in the Store.")
	}

	bsecret, _ := store.Secret(uid)

	if asecret != bsecret {
		t.Error("Mismatch between the secret from the same user.")
	}
}

func TestUid(t *testing.T) {
	store := mock.NewTestStore()

	userA := store.CreateUser()

	userB := store.CreateUser()

	if userA == userB {
		t.Error("Two users cannot have the same id.")
	}
}

func TestMacaroon(t *testing.T) {
	store := mock.NewTestStore()

	uid := store.CreateUser()

	secret, _ := store.Secret(uid)

	oven := macaroon.NewOven(secret)

	mac, _ := oven.MapCaveats(caveats).Service(service).Cook()

	signaturea := mac.Signature()

	uid = store.CreateUser()

	secret, _ = store.Secret(uid)

	oven = macaroon.NewOven(secret)

	mac, _ = oven.MapCaveats(caveats).Service(service).Cook()

	signatureb := mac.Signature()

	if signaturea == signatureb {
		t.Error("Two users cannot produce the same signature for the same given caveats.")
	}
}
