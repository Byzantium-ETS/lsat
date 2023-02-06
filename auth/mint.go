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

// Assure qu'un macaroon ait accès aux services
type ServiveHandler interface {
	Capabilities(context.Context, Service) ([]Caveat, error) // The capabilities of the service
}

// https://github.com/lightninglabs/aperture/blob/master/mint/mint.go#L65
type Minter struct {
	// Une abstraction des services offert par une application
	service ServiveHandler

	// La source des secrets des lsats qui seront créé
	secrets secrets.SecretStore

	// Crée les nouveaux invoices
	// qui serviront de challenge pour créer les lsat
	node lightning.LightningNode
}
