package tests

import (
	"lsat/mock"
	"lsat/secrets"
	"testing"
)

var secretStore mock.TestStore = mock.NewTestStore()

func TestMakeSecret(t *testing.T) {
	root := secrets.NewSecret()

	t.Log(len(root[:]))
	t.Log(root)

	newRoot, err := secrets.MakeSecret(root[:])

	if err != nil {
		t.Error(err)
	}

	if newRoot != root {
		t.Log("both secrets should be equal")
	}

}

func TestSecret(t *testing.T) {
	uid := secretStore.CreateUser()

	asecret, err := secretStore.NewSecret(uid)

	if err != nil {
		t.Error("Secret not found in the Store.")
	}

	bsecret, _ := secretStore.GetSecret(uid)

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
