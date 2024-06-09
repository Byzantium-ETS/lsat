package tests

import (
	"lsat/secrets"
	"testing"
)

var secretStore = secrets.NewSecretFactory()

func TestGetSecret(t *testing.T) {
	user := secretStore.NewUser()
	secretA, err := secretStore.NewSecret(user)

	if err != nil {
		t.Error(err)
	}

	secretB, err := secretStore.GetSecret(user)

	if err != nil {
		t.Error(err)
	}

	if secretA != secretB {
		t.Error("The secret retrieved from the user should be same created.")
	}

}

func TestNewSecret(t *testing.T) {
	userA := secretStore.NewUser()

	userB := secretStore.NewUser()

	secretA, err := secretStore.NewSecret(userA)

	if err != nil {
		t.Error(err)
	}

	secretB, err := secretStore.NewSecret(userB)

	if err != nil {
		t.Error(err)
	}

	if secretA == secretB {
		t.Error("Two users cannot have the same secret")
	}
}

func TestNewUser(t *testing.T) {
	userA := secretStore.NewUser()

	userB := secretStore.NewUser()

	if userA == userB {
		t.Error("Two users cannot have the same id.")
	}
}
