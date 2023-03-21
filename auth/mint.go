package auth

import (
	"context"
	"errors"
	"lsat/lightning"
	"lsat/macaroon"
	"lsat/secrets"
)

const (
	PermErr  = "The macaroon lacks permission!"
	TokenErr = "The token could not be found!"
	UserErr  = "The macaroon is not from this user!"
)

// Assure qu'un macaroon ait accès aux services
type ServiveManager interface {
	Services(cx context.Context, names ...string) ([]macaroon.Service, error)
	Capabilities(cx context.Context, services ...macaroon.Service) ([]macaroon.Caveat, error) // The capabilities of the service
	VerifyCaveats(cx macaroon.Service, caveats ...macaroon.Caveat) error                      // Retourne une erreur si les caveats ne sont plus valide, ca veut generalement dire que le macaroon est expire
}

type SecretStore interface {
	// NewSecret() (Secret, error)
	Secret(uid secrets.UserId) (secrets.Secret, error)          // S'il n'y a pas de RootKey pour l'utilisateur, il sera créé
	StoreToken(id macaroon.TokenID, token macaroon.Token) error // Les tokens peuvent être conservé pour des raisons d'archivage
	Tokens() *map[macaroon.TokenID]macaroon.Token               // Tous les tokens d'un utilisateur
}

// https://github.com/lightninglabs/aperture/blob/master/mint/mint.go#L65
type Minter struct {
	service    ServiveManager       // Une abstraction des services offert par une application
	secrets    SecretStore          // La source des secrets des lsats qui seront créé
	challenger lightning.Challenger // Crée les challenges sous la forme d'invoices
}

func NewMinter(service ServiveManager, secrets SecretStore, challenger lightning.Challenger) Minter {
	return Minter{service, secrets, challenger}
}

func (minter *Minter) SecretStore() SecretStore {
	return minter.secrets
}

func totalPrice(services ...macaroon.Service) int {
	total := 0
	for _, s := range services {
		total += int(s.Price)
	}
	return total
}

func (minter *Minter) MintToken(uid secrets.UserId, service_names ...string) (macaroon.PreToken, error) {
	token := macaroon.PreToken{}

	services, err := minter.service.Services(context.Background(), service_names...)

	if err != nil {
		return token, err
	}

	preimage, payment, err := minter.challenger.Challenge(int64(totalPrice(services...)))

	if err != nil {
		return token, err
	}

	token.Invoice = payment.Invoice

	caveats, err := minter.service.Capabilities(context.Background(), services...)

	secret, err := minter.secrets.Secret(uid)

	if err != nil {
		return token, err
	}

	oven := macaroon.NewOven(secret)

	mac, _ := oven.MapCaveats(caveats).Cook()

	token.Mac = mac

	lsat := macaroon.Token{Mac: mac, Preimage: preimage}

	tokenId := macaroon.NewTokenID(uid, preimage.Hash())

	minter.secrets.StoreToken(tokenId, lsat)

	return token, nil
}

func (minter *Minter) authToken(uid secrets.UserId, lsat *macaroon.Token) error {
	tokens := minter.secrets.Tokens()

	tokenId := macaroon.NewTokenID(uid, lsat.Preimage.Hash())

	rlsat, ok := (*tokens)[tokenId]

	if !ok || &rlsat != lsat {
		return errors.New(TokenErr)
	}

	return minter.verifyMacaroon(uid, &lsat.Mac)
}

func (minter *Minter) verifyMacaroon(uid secrets.UserId, mac *macaroon.Macaroon) error {
	secret, _ := minter.secrets.Secret(uid)
	oven := macaroon.NewOven(secret)
	nmac, _ := oven.MapCaveats(mac.Caveats()).Cook()
	if mac.Signature() == nmac.Signature() {
		return minter.service.VerifyCaveats(mac.Service(), mac.Caveats()...)
	} else {
		return errors.New(UserErr) // Faudrait des erreurs
	}
}
