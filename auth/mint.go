package auth

import (
	"context"
	"errors"
	"lsat/lightning"
	"lsat/secrets"
)

type Ressource interface {
}

// Assure qu'un macaroon ait accès aux services
type ServiveManager interface {
	Services(context.Context, ...string) []Service
	Capabilities(context.Context, ...Service) ([]Caveat, error) // The capabilities of the service
}

// https://github.com/lightninglabs/aperture/blob/master/mint/mint.go#L65
type Minter struct {
	// Une abstraction des services offert par une application
	service ServiveManager

	// La source des secrets des lsats qui seront créé
	secrets secrets.SecretStore

	// Crée les challenges sous la forme d'invoices
	challenger lightning.Challenger
}

func totalPrice(services ...Service) int {
	total := 0
	for _, s := range services {
		total += int(s.Price)
	}
	return total
}

func (minter *Minter) GetRessource(uid secrets.UserId, token LSAT) (Ressource, error) {
	return recover(), nil
}

func (minter *Minter) MintLSAT(uid secrets.UserId, services ...Service) (preLSAT, error) {
	// Enregistrer le LSAT valide dans le SecretStore
	// minter.secrets.StoreLSAT(uid, LSAT{})
	return preLSAT{}, nil
}

func (minter *Minter) verifyLSAT(uid secrets.UserId, token LSAT) error {
	oven, _ := NewOven(&minter.secrets, uid)
	oven.MapCaveats(token.Mac.caveats)
	_, mac, _ := oven.Cook()
	if mac.Signature() == token.Mac.Signature() {
		return nil
	} else {
		return errors.New("Invalid macaroon") // Faudrait des erreurs
	}
}
