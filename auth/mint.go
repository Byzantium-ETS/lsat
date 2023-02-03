package auth

import (
	"lsat/secrets"
)

type Minter struct {
	store secrets.RootKeyStore
}
