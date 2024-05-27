package macaroon

import (
	"crypto/hmac"
	"crypto/sha256"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lntypes"
)

// Oven bakes macaroons by combining the root secret, user ID, and caveats.
type Oven struct {
	userId   secrets.UserId
	root     secrets.Secret
	caveats  []Caveat
	macaroon *Macaroon
}

// Creates a new Oven with the given root secret.
func NewOven(root secrets.Secret) Oven {
	oven := Oven{}
	oven.root = root
	return oven
}

// Sets the user ID in the Oven.
func (oven Oven) WithUserId(uid secrets.UserId) Oven {
	oven.userId = uid
	return oven
}

// Adds a single caveat to the Oven.
func (oven *Oven) Attenuate(caveat Caveat) {
	oven.caveats = append(oven.caveats, caveat)
}

// Adds multiple third party caveats to the Oven.
func (oven Oven) WithThirdPartyCaveats(caveats ...Caveat) Oven {
	oven.caveats = append(oven.caveats, caveats...)
	return oven
}

// Adds multiple first party caveats to the Oven.
func (oven Oven) WithFirstPartyCaveats(caveats ...Caveat) Oven {
	oven.caveats = append(caveats, oven.caveats...)
	return oven
}

// Adds caveats for the specified services to the Oven.
func (oven Oven) WithService(services ...Service) Oven {
	for _, service := range services {
		oven.caveats = append([]Caveat{service.Caveat()}, oven.caveats...)
	}
	return oven
}

// Cook computes the signature of the Macaroon and returns it.
func (oven Oven) Cook() (Macaroon, error) {
	// Create a new HMAC with SHA-256 using the root secret as the key
	mac := hmac.New(sha256.New, oven.root[:])

	// Write the string representation of each caveat into the HMAC
	for i, caveat := range oven.caveats {
		if i > 0 {
			mac = hmac.New(sha256.New, mac.Sum(nil))
		}
		mac.Write([]byte(caveat.String()))
	}

	// Generate a signature by summing the HMAC
	signature, err := lntypes.MakeHash(mac.Sum(nil))

	if err != nil {
		return Macaroon{}, err
	}

	var caveats []Caveat

	if oven.macaroon != nil {
		caveats = append(oven.macaroon.caveats, oven.caveats...)
	} else {
		caveats = oven.caveats
	}

	// Create and return the Macaroon with the user ID, caveats, and signature
	return Macaroon{userId: oven.userId, caveats: caveats, signature: signature}, nil
}
