package auth

import (
	"context"
	"lsat/lightning"
	"lsat/secrets"
)

type Service struct {
	Name string

	Price int64
}

type Ressource interface {
}

// Assure qu'un macaroon ait accès aux services
type ServiveManager interface {
	Capabilities(context.Context, Service) ([]Caveat, error) // The capabilities of the service
}

// https://github.com/lightninglabs/aperture/blob/master/mint/mint.go#L65
type Minter struct {
	// Une abstraction des services offert par une application
	service ServiveManager

	// La source des secrets des lsats qui seront créé
	secrets secrets.SecretStore

	// Crée les nouveaux invoices
	// qui serviront de challenge pour créer les lsat
	node lightning.LightningNode
}

func (minter Minter) mintToken(uid secrets.UserId, caveat []Caveat) (preToken, error) {
	return preToken{}, nil
}
