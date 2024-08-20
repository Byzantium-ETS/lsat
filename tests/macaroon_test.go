package tests

import (
	"lsat/macaroon"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var caveat macaroon.Caveat = macaroon.NewCaveat(macaroon.ExpiryDateKey, time.Now().Add(time.Hour).Format(time.RFC3339))

func TestSignature(t *testing.T) {
	uid := secretStore.NewUser()

	secret, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(secret)

	mac, _ := oven.WithThirdPartyCaveats(caveat).Cook()

	signaturea := mac.Signature()

	uid = secretStore.NewUser()

	secret, _ = secretStore.NewSecret(uid)

	oven = macaroon.NewOven(secret)

	mac, _ = oven.WithThirdPartyCaveats(caveat).Cook()

	signatureb := mac.Signature()

	if signaturea == signatureb {
		t.Error("Two users cannot produce the same signature for the same given caveats.")
	}
}

func TestValueIter(t *testing.T) {
	uid := secretStore.NewUser()

	secret, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(secret)

	mac, _ := oven.WithUserId(uid).WithThirdPartyCaveats(macaroon.NewCaveat("name", "bob")).WithThirdPartyCaveats(caveat).Cook()

	iter := mac.GetValue(macaroon.ExpiryDateKey)
	expiryDate := iter.Next()

	t.Log(expiryDate)

	if expiryDate != caveat.Value {
		t.Error("The expiry date is not the same as the one in the caveat.")
	}
}

func TestMacaroonEncoding(t *testing.T) {
	uid := secretStore.NewUser()
	root, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(root)
	oven = oven.WithUserId(uid).WithThirdPartyCaveats(macaroon.NewCaveat("name", "bob"))

	mac, err := oven.Cook()

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(mac.ToJSON())

	encodedMac := mac.String()

	t.Log(encodedMac)

	decodedMac, err := macaroon.DecodeBase64(encodedMac)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, mac, decodedMac)
}

func TestMacaroonSignature(t *testing.T) {
	uid := secretStore.NewUser()
	root, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(root)
	oven = oven.WithUserId(uid).WithThirdPartyCaveats(macaroon.NewCaveat("name", "bob"))

	mac1, err := oven.Cook()

	if err != nil {
		t.Error(err)
		return
	}

	uid = secretStore.NewUser()
	root, _ = secretStore.NewSecret(uid)

	oven = macaroon.NewOven(root)
	oven = oven.WithUserId(uid).WithThirdPartyCaveats(macaroon.NewCaveat("name", "bob"))

	mac2, err := oven.Cook()

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(mac1)
	t.Log(mac2)

	assert.NotEqual(t, mac1.Signature(), mac2.Signature())
}

func TestFirstPartyCaveats(t *testing.T) {
	uid := secretStore.NewUser()
	root, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(root)
	oven = oven.WithUserId(uid).WithThirdPartyCaveats(macaroon.NewCaveat("name", "bob"))

	mac1, _ := oven.WithFirstPartyCaveats(testService.Caveats()...).Cook()
	mac2, _ := oven.WithFirstPartyCaveats(testService.Caveats()...).Cook()

	t.Log(mac1.ToJSON())
	t.Log(mac2.ToJSON())

	assert.Equal(t, mac1, mac2)
}

func TestThirdPartyCaveats(t *testing.T) {
	uid := secretStore.NewUser()
	root, _ := secretStore.NewSecret(uid)

	oven := macaroon.NewOven(root)
	oven = oven.WithUserId(uid).WithThirdPartyCaveats(macaroon.NewCaveat("name", "bob"))

	mac, _ := oven.Cook()

	t.Log(len(mac.Signature()))

	thirdPartyCaveat := macaroon.NewCaveat("color", "red")

	macThirdParty, _ := mac.Oven().WithThirdPartyCaveats(thirdPartyCaveat).Cook()
	macFirstParty, _ := oven.WithThirdPartyCaveats(thirdPartyCaveat).Cook()

	t.Log(macFirstParty.ToJSON())
	t.Log(macThirdParty.ToJSON())

	assert.Equal(t, macFirstParty, macThirdParty)
}
