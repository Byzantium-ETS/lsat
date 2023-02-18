package auth

import (
	"context"
	"errors"
	"lsat/lightning"
	"lsat/secrets"
)

const (
	PermErr  = "The macaroon lacks permission!"
	TokenErr = "The token could not be found!"
	UserErr  = "The macaroon is not from this user!"
)

// Assure qu'un macaroon ait accès aux services
type ServiveManager interface {
	Services(context.Context, ...string) ([]Service, error)
	Capabilities(context.Context, ...Service) ([]Caveat, error) // The capabilities of the service
	VerifyCaveats(Service, ...Caveat) error                     // Retourne une erreur si les caveats ne sont plus valide, ca veut generalement dire que le macaroon est expire
}

// https://github.com/lightninglabs/aperture/blob/master/mint/mint.go#L65
type Minter struct {
	service    ServiveManager       // Une abstraction des services offert par une application
	secrets    secrets.SecretStore  // La source des secrets des lsats qui seront créé
	challenger lightning.Challenger // Crée les challenges sous la forme d'invoices
}

func totalPrice(services ...Service) int {
	total := 0
	for _, s := range services {
		total += int(s.Price)
	}
	return total
}

func (minter *Minter) MintToken(uid secrets.UserId, service_names ...string) (preToken, error) {
	token := preToken{}

	services, _ := minter.service.Services(context.Background(), service_names...)

	preimage, payment, err := minter.challenger.Challenge(int64(totalPrice(services...)))

	if err != nil {
		return token, err
	}

	token.Invoice = payment.Invoice

	caveats, err := minter.service.Capabilities(context.Background(), services...)

	oven, _ := NewOven(&minter.secrets, uid)

	mac, _ := oven.MapCaveats(caveats).Cook()

	token.Mac = mac

	lsat := Token{Mac: mac, Preimage: preimage}

	tokenId := NewTokenID(uid, preimage.Hash())

	minter.secrets.StoreToken(tokenId, lsat)

	return token, nil
}

func (minter *Minter) authToken(uid secrets.UserId, lsat *Token) error {
	tokens := minter.secrets.Tokens()

	tokenId := NewTokenID(uid, lsat.Preimage.Hash())

	rlsat, ok := tokens[tokenId]

	if !ok || rlsat != lsat {
		return errors.New(TokenErr)
	}

	return minter.verifyMacaroon(uid, &lsat.Mac)
}

func (minter *Minter) verifyMacaroon(uid secrets.UserId, mac *Macaroon) error {
	oven, _ := NewOven(&minter.secrets, uid)
	oven.MapCaveats(mac.caveats)
	nmac, _ := oven.Cook()
	if mac.Signature() == nmac.Signature() {
		return minter.service.VerifyCaveats(mac.service, mac.caveats...)
	} else {
		return errors.New(UserErr) // Faudrait des erreurs
	}
}
