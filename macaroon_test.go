package main

import (
	"lsat/macaroon"
	"reflect"
	"testing"
)

func TestMacaroonEncoding(t *testing.T) {
	uid := secretStore.CreateUser()
	root, _ := secretStore.Secret(uid)

	oven := macaroon.NewOven(root)
	oven = oven.UserId(uid).Caveats(macaroon.NewCaveat("name", "bob")).Service(macaroon.NewService("rent", 1000))

	mac, err := oven.Cook()

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(mac.Caveats())

	encodedMac := mac.String()

	t.Log(encodedMac)

	decodedMac, err := macaroon.Decode(encodedMac)

	if err != nil {
		t.Error(err)
	} else if reflect.DeepEqual(decodedMac, mac) {
		t.Error("Failed to decode macaroon!")
	}
}

func TestMacaroonSignature(t *testing.T) {
	uid := secretStore.CreateUser()
	root, _ := secretStore.Secret(uid)

	oven := macaroon.NewOven(root)
	oven = oven.UserId(uid).Caveats(macaroon.NewCaveat("name", "bob")).Service(macaroon.NewService("rent", 1000))

	mac1, err := oven.Cook()

	if err != nil {
		t.Error(err)
		return
	}

	uid = secretStore.CreateUser()
	root, _ = secretStore.Secret(uid)

	oven = macaroon.NewOven(root)
	oven = oven.UserId(uid).Caveats(macaroon.NewCaveat("name", "bob")).Service(macaroon.NewService("rent", 1000))

	mac2, err := oven.Cook()

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(mac1)
	t.Log(mac2)

	if mac1.String() == mac2.String() {
		t.Error("The hex encoding cannot be similar!")
	}

}
