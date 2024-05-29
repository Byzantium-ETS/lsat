package tests

import (
	"lsat/auth"
	"lsat/macaroon"
	"lsat/secrets"
	"testing"

	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/stretchr/testify/assert"
)

func TestStoreToken(t *testing.T) {
	store, err := auth.NewStore("../.store")

	if err != nil {
		t.Error(err)
	}

	secret := secrets.NewSecret()
	preimage, _ := lntypes.MakePreimage(secret[:])

	token := macaroon.Token{
		Macaroon: macaroon.Macaroon{},
		Preimage: preimage,
	}

	id := token.Id()

	err = store.StoreToken(id, token)

	if err != nil {
		t.Error(err)
	}

	store.RemoveToken(id)
}

func TestGetToken(t *testing.T) {
	store, err := auth.NewStore("../.store")

	if err != nil {
		t.Error(err)
	}

	secret := secrets.NewSecret()
	preimage, _ := lntypes.MakePreimage(secret[:])

	tokenIn := macaroon.Token{
		Macaroon: macaroon.Macaroon{},
		Preimage: preimage,
	}

	t.Log(tokenIn)

	id := tokenIn.Id()

	_ = store.StoreToken(id, tokenIn)

	tokenOut, err := store.GetToken(id)

	if err != nil {
		t.Error(err)
	}

	t.Log(tokenOut)

	assert.Equal(t, tokenIn.String(), tokenOut.String())

	store.RemoveToken(id)
}
