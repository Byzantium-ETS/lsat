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
	PermErr  = "The macaroon lacks permissions!"
	TokenErr = "The token could not be found!"
	SigErr   = "The macaroon has an invalid signature!"
)

// https://github.com/lightninglabs/aperture/blob/master/mint/mint.go#L65
type Minter struct {
	service    ServiceLimiter       // Une abstraction des services offert par une application
	secrets    SecretStore          // La source des secrets des lsats qui seront créé
	challenger challenge.Challenger // Crée les challenges sous la forme d'invoices
}

func NewMinter(service ServiceLimiter, secrets SecretStore, challenger challenge.Challenger) Minter {
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

	token.PaymentRequest = result.PaymentRequest

	caveats, err := minter.service.Capabilities(context.Background(), services...)

	secret, err := minter.secrets.GetSecret(uid)

	if err != nil {
		secret, _ = minter.secrets.NewSecret(uid)
	}

	if err != nil {
		fmt.Println(err)
	}

	oven := macaroon.NewOven(secret)

	mac, err := oven.WithUserId(uid).WithService(services...).WithCaveats(caveats...).Cook()

	if err != nil {
		return token, err
	}

	token.Macaroon = mac

	lsat := macaroon.Token{Macaroon: mac, Preimage: result.Preimage}

	tokenId := lsat.Id()

	// We store the token in an archive.
	minter.secrets.StoreToken(tokenId, lsat)

	return token, nil
}

func (minter *Minter) AuthToken(lsat *macaroon.Token) error {
	tokens := *minter.secrets.Tokens()

	_, ok := tokens[lsat.Id()]

	if !ok {
		return errors.New(TokenErr)
	}

	return minter.AuthMacaroon(&lsat.Macaroon)
}

func (minter *Minter) AuthMacaroon(mac *macaroon.Macaroon) error {
	secret, _ := minter.secrets.GetSecret(mac.UserId())
	oven := macaroon.NewOven(secret)
	nmac, _ := oven.WithCaveats(mac.Caveats()...).Cook()
	if mac.Signature() == nmac.Signature() {
		return minter.service.VerifyCaveats(mac.Caveats()...)
	} else {
		return errors.New(SigErr) // Faudrait des erreurs
	}
}
