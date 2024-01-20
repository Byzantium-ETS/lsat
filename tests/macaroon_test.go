package tests

import (
	"lsat/macaroon"
	"reflect"
	"testing"
)

func TestMacaroon(t *testing.T) {
	uid := secretStore.CreateUser()

	secret, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(secret)

	mac, _ := oven.WithCaveats(caveat).WithService(service).Cook()

	signaturea := mac.Signature()

	uid = secretStore.CreateUser()

	secret, _ = secretStore.NewSecret(uid)

	oven = macaroon.NewOven(secret)

	mac, _ = oven.WithCaveats(caveat).WithService(service).Cook()

	signatureb := mac.Signature()

	if signaturea == signatureb {
		t.Error("Two users cannot produce the same signature for the same given caveats.")
	}
}

func TestMacaroonEncoding(t *testing.T) {
	uid := secretStore.CreateUser()
	root, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(root)
	oven = oven.WithUserId(uid).WithCaveats(macaroon.NewCaveat("name", "bob")).WithService(macaroon.NewService("rent", 1000))

	mac, err := oven.Cook()

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(mac.ToJSON())

	encodedMac := mac.ToJSON().String()

	t.Log(encodedMac)

	decodedMac, err := macaroon.DecodeBase64(encodedMac)

	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(decodedMac, mac) {
		t.Log(decodedMac.ToJSON())
		t.Error("Failed to decode macaroon!")
	}
}

func TestMacaroonSignature(t *testing.T) {
	uid := secretStore.CreateUser()
	root, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(root)
	oven = oven.WithUserId(uid).WithCaveats(macaroon.NewCaveat("name", "bob")).WithService(macaroon.NewService("rent", 1000))

	mac1, err := oven.Cook()

	if err != nil {
		t.Error(err)
		return
	}

	uid = secretStore.CreateUser()
	root, _ = secretStore.NewSecret(uid)

	oven = macaroon.NewOven(root)
	oven = oven.WithUserId(uid).WithCaveats(macaroon.NewCaveat("name", "bob")).WithService(macaroon.NewService("rent", 1000))

	mac2, err := oven.Cook()

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(mac1)
	t.Log(mac2)

	if mac1.Signature() == mac2.Signature() {
		t.Error("The hex encoding cannot be similar!")
	}

}

func TestCaveats(t *testing.T) {
	uid := secretStore.CreateUser()
	root, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(root)
	oven = oven.WithUserId(uid).WithService(macaroon.NewService("rent", 1000)).WithCaveats(macaroon.NewCaveat("name", "bob"))

	caveat := macaroon.NewCaveat("color", "red")

	mac1, _ := oven.WithCaveats(caveat).Cook()
	mac2, _ := oven.WithCaveats(caveat).Cook()

	t.Log(mac1.ToJSON())
	t.Log(mac2.ToJSON())

	if !reflect.DeepEqual(mac1, mac2) {
		t.Error("Both macaroons should have the same signature!")
	}
}

func TestThirdPartyCaveats(t *testing.T) {
	uid := secretStore.CreateUser()
	root, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(root)
	oven = oven.WithUserId(uid).WithService(macaroon.NewService("rent", 1000)).WithCaveats(macaroon.NewCaveat("name", "bob"))

	mac, _ := oven.Cook()

	t.Log(len(mac.Signature()))

	thirdPartyCaveat := macaroon.NewCaveat("color", "red")

	macThirdParty, _ := mac.Oven().WithCaveats(thirdPartyCaveat).Cook()
	macFirstParty, _ := oven.WithCaveats(thirdPartyCaveat).Cook()

	t.Log(macFirstParty.ToJSON())
	t.Log(macThirdParty.ToJSON())

	if !reflect.DeepEqual(macFirstParty, macThirdParty) {
		t.Error("Both macaroons should have the same signature!")
	}
}
