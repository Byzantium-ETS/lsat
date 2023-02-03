package secrets

import (
	"lsat/auth"
)

type RootKey struct {
	root []int8
}

type RootKeyStore interface {
	GetRoot(uid auth.MacaroonId) (RootKey, error)
	CreateRoot(uid auth.MacaroonId) (RootKey, error)
}
