package auth

import (
	"context"
	"errors"
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
type ServiveLimiter interface {
	Capabilities(context.Context, ...Service) ([]Caveat, error) // The capabilities of the service
}

// https://github.com/lightninglabs/aperture/blob/master/mint/mint.go#L65
type Minter struct {
	// Une abstraction des services offert par une application
	service ServiveLimiter

	// La source des secrets des lsats qui seront créé
	secrets secrets.SecretStore

	// Crée les nouveaux invoices
	// qui serviront de challenge pour créer les lsat
	node lightning.LightningNode
}

func (minter *Minter) GetRessource(uid secrets.UserId, token Token) (Ressource, error) {
	return recover(), nil
}

func (minter *Minter) MintToken(uid secrets.UserId, services ...Service) (preToken, error) {
	return preToken{}, nil
}

func (minter *Minter) verifyToken(uid secrets.UserId, token Token) error {
	oven, _ := NewOven(&minter.secrets, uid)
	oven.MapCaveats(token.Mac.caveats)
	mac, _ := oven.Cook()
	if mac.Signature() == token.Mac.Signature() {
		return nil
	} else {
		return errors.New("Invalid macaroon") // Faudrait des erreurs
	}
}
