package auth

import (
	"errors"
	"lsat/challenge"
	"lsat/macaroon"
	"lsat/secrets"
	"lsat/service"
)

const (
	permErr = "the macaroon lacks permissions"
	hashErr = "the payment_hash does not correspond to the preimage"
	sigErr  = "the macaroon has an invalid signature"
)

type Minter struct {
	service    service.ServiceManager
	secrets    secrets.SecretStore
	challenger challenge.Challenger
}

func NewMinter(service service.ServiceManager, secrets secrets.SecretStore, challenger challenge.Challenger) Minter {
	return Minter{service, secrets, challenger}
}

func (minter *Minter) SecretStore() secrets.SecretStore {
	return minter.secrets
}

func totalPrice(services ...service.Service) uint64 {
	var total uint64 = 0
	for _, s := range services {
		total += s.Price
	}
	return total
}

// MintToken generates a new pre-token for the user.
func (minter *Minter) MintToken(uid secrets.UserId, service_id service.ServiceId) (macaroon.PreToken, error) {
	// Initialize an empty pre-token.
	token := macaroon.PreToken{}

	// Fetch information about the requested services.
	service, err := minter.service.Service(service_id)
	if err != nil {
		return token, err
	}

	// Initiate a payment challenge using the total price of the requested services.
	result, err := minter.challenger.Challenge(totalPrice(service))
	if err != nil {
		return token, err
	}

	// Set the PaymentRequest in the pre-token based on the result of the payment challenge.
	token.InvoiceResponse = result

	// Retrieve the capabilities (caveats) associated with the requested services.
	caveats := service.Caveats()

	// Get or create a secret associated with the user ID.
	secret, err := minter.secrets.GetSecret(uid)
	if err != nil {
		// If an error occurs, create a new secret for the user ID.
		secret, _ = minter.secrets.NewSecret(uid)
	}

	// Create a Macaroon oven with the obtained or created secret.
	oven := macaroon.NewOven(secret)

	// Cook the Macaroon with the user ID, requested services, and retrieved capabilities.
	mac, err := oven.WithUserId(uid).WithFirstPartyCaveats(caveats...).Cook()
	if err != nil {
		return token, err
	}

	// Store the Macaroon in the secrets archive.
	token.Macaroon, _ = mac.Oven().WithFirstPartyCaveats(macaroon.Caveat{
		Key: macaroon.PaymentHashKey, Value: result.PaymentHash.String(),
	}).Cook()

	// Return the generated pre-token.
	return token, nil
}

// Authentify the validity of the token.
func (minter *Minter) AuthToken(lsat *macaroon.Token) error {
	// Verify the preimage
	paymentHash := lsat.Macaroon.GetValue(macaroon.PaymentHashKey)
	if len(paymentHash) == 0 {
		return errors.New(permErr)
	}

	if lsat.Preimage.Hash().String() != paymentHash[0] {
		return errors.New(hashErr)
	}

	// Validate the LSAT's Macaroon using the authentication service.
	err := minter.AuthMacaroon(&lsat.Macaroon)
	if err != nil {
		return err
	}

	return nil
}

// Verifies that signature and caveats are valid.
func (minter *Minter) AuthMacaroon(mac *macaroon.Macaroon) error {
	secret, _ := minter.secrets.GetSecret(mac.UserId())
	oven := macaroon.NewOven(secret)
	nmac, _ := oven.WithThirdPartyCaveats(mac.Caveats()...).Cook()

	if mac.Signature() != nmac.Signature() {
		return errors.New(sigErr) // Faudrait des erreurs
	}

	err := minter.service.VerifyCaveats(mac.Caveats()...) /// Check the validity of the caveats

	if err != nil {
		return err
	}

	return nil
}
