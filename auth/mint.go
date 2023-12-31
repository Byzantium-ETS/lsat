package auth

import (
	"context"
	"errors"
	"fmt"
	"lsat/challenge"
	"lsat/macaroon"
	"lsat/secrets"
)

const (
	PermErr    = "The macaroon lacks permission!"
	TokenErr   = "The token could not be found!"
	SigErr     = "The macaroon has an invalid signature!"
	PaymentErr = "The preimage provided is invalid!"
)

// Assure qu'un macaroon ait accès aux services
type ServiveManager interface {
	Services(cx context.Context, names ...string) ([]macaroon.Service, error)
	Capabilities(cx context.Context, services ...macaroon.Service) ([]macaroon.Caveat, error) // The capabilities of the service
	VerifyCaveats(caveats ...macaroon.Caveat) error                                           // Retourne une erreur si les caveats ne sont plus valide, ca veut generalement dire que le macaroon est expire
}

type SecretStore interface {
	Secret(uid secrets.UserId) (secrets.Secret, error)          // S'il n'y a pas de RootKey pour l'utilisateur, il sera créé
	StoreToken(id macaroon.TokenID, token macaroon.Token) error // Les tokens peuvent être conservé pour des raisons d'archivage
	Tokens() *map[macaroon.TokenID]macaroon.Token               // Tous les tokens d'un utilisateur
}

// https://github.com/lightninglabs/aperture/blob/master/mint/mint.go#L65
type Minter struct {
	service    ServiveManager       // Une abstraction des services offert par une application
	secrets    SecretStore          // La source des secrets des lsats qui seront créé
	challenger challenge.Challenger // Crée les challenges sous la forme d'invoices
}

func NewMinter(service ServiveManager, secrets SecretStore, challenger challenge.Challenger) Minter {
	return Minter{service, secrets, challenger}
}

func (minter *Minter) SecretStore() SecretStore {
	return minter.secrets
}

func totalPrice(services ...macaroon.Service) uint64 {
	var total uint64 = 0
	for _, s := range services {
		total += s.Price
	}
	return total
}

func (minter *Minter) MintToken(uid secrets.UserId, service_names ...string) (macaroon.PreToken, error) {
	token := macaroon.PreToken{}

	services, err := minter.service.Services(context.Background(), service_names...)

	if err != nil {
		return token, err
	}

	result, err := minter.challenger.Challenge(totalPrice(services...))

	if err != nil {
		return token, err
	}

	token.Invoice = result.Invoice

	caveats, err := minter.service.Capabilities(context.Background(), services...)

	secret, err := minter.secrets.Secret(uid)

	if err != nil {
		fmt.Println(err)
	}

	oven := macaroon.NewOven(secret)

	mac, err := oven.UserId(uid).Caveats(caveats...).Service(services...).Cook()

	if err != nil {
		return token, err
	}

	token.Mac = mac

	lsat := macaroon.Token{Mac: mac, Preimage: result.Preimage}

	tokenId := lsat.Id()

	// We store the token in an archive.
	minter.secrets.StoreToken(tokenId, lsat)

	return token, nil
}

func (minter *Minter) AuthToken(lsat *macaroon.Token) error {
	tokens := *minter.secrets.Tokens()

	_, ok := tokens[lsat.Id()]

	if !ok {
		return errors.New(PaymentErr)
	}

	return minter.AuthMacaroon(&lsat.Mac)
}

func (minter *Minter) AuthMacaroon(mac *macaroon.Macaroon) error {
	secret, _ := minter.secrets.Secret(mac.Uid())
	oven := macaroon.NewOven(secret)
	nmac, _ := oven.Caveats(mac.Caveats()...).Cook()
	if mac.Signature() == nmac.Signature() {
		return minter.service.VerifyCaveats(mac.Caveats()...)
	} else {
		return errors.New(SigErr) // Faudrait des erreurs
	}
}
