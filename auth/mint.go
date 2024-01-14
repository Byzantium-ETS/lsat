package auth

import (
	"context"
	"errors"
	"lsat/challenge"
	"lsat/macaroon"
	"lsat/secrets"
)

const (
	permErr  = "The macaroon lacks permissions!"
	tokenErr = "The token could not be found!"
	sigErr   = "The macaroon has an invalid signature!"
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

// MintToken generates a new pre-token for the user.
func (minter *Minter) MintToken(uid secrets.UserId, service_names ...string) (macaroon.PreToken, error) {
	// Initialize an empty pre-token.
	token := macaroon.PreToken{}

	// Fetch information about the requested services.
	services, err := minter.service.Services(context.Background(), service_names...)
	if err != nil {
		return token, err
	}

	// Initiate a payment challenge using the total price of the requested services.
	result, err := minter.challenger.Challenge(totalPrice(services...))
	if err != nil {
		return token, err
	}

	// Set the PaymentRequest in the pre-token based on the result of the payment challenge.
	token.PaymentRequest = result.PaymentRequest

	// Retrieve the capabilities (caveats) associated with the requested services.
	caveats, err := minter.service.Capabilities(context.Background(), services...)
	if err != nil {
		return token, err
	}

	// Get or create a secret associated with the user ID.
	secret, err := minter.secrets.GetSecret(uid)
	if err != nil {
		// If an error occurs, create a new secret for the user ID.
		secret, _ = minter.secrets.NewSecret(uid)
	}

	// Create a Macaroon oven with the obtained or created secret.
	oven := macaroon.NewOven(secret)

	// Cook the Macaroon with the user ID, requested services, and retrieved capabilities.
	mac, err := oven.WithUserId(uid).WithService(services...).WithCaveats(caveats...).Cook()
	if err != nil {
		return token, err
	}

	// Store the Macaroon in the secrets archive.
	token.Macaroon = mac

	// Create an LSAT with the Macaroon and the preimage obtained from the payment challenge.
	lsat := macaroon.Token{Macaroon: mac, Preimage: result.Preimage}

	// Store the LSAT in the secrets archive.
	tokenId := lsat.Id()
	minter.secrets.StoreToken(tokenId, lsat)

	// Return the generated pre-token.
	return token, nil
}

// AuthToken generates an authentication token (Macaroon) based on the provided LSAT (Lightning Service Authentication Token).
func (minter *Minter) AuthToken(lsat *macaroon.Token) (macaroon.Macaroon, error) {
	// Retrieve the stored tokens.
	tokens := *minter.secrets.Tokens()

	// Check if the LSAT's ID is present in the stored tokens.
	_, ok := tokens[lsat.Id()]
	if !ok {
		return macaroon.Macaroon{}, errors.New(tokenErr)
	}

	// Validate the LSAT's Macaroon using the authentication service.
	err := minter.AuthMacaroon(&lsat.Macaroon)
	if err != nil {
		return macaroon.Macaroon{}, err
	}

	// Sign the validated Macaroon using the service's Sign method.
	return minter.service.Sign(lsat.Macaroon)
}

// / AuthMacaroon only verifies that the authorization server has minter the macaroon.
func (minter *Minter) AuthMacaroon(mac *macaroon.Macaroon) error {
	secret, _ := minter.secrets.GetSecret(mac.UserId())
	oven := macaroon.NewOven(secret)
	nmac, _ := oven.WithCaveats(mac.Caveats()...).Cook()

	if mac.Signature() != nmac.Signature() {
		return errors.New(sigErr) // Faudrait des erreurs
	}

	return nil
}
