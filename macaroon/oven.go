package macaroon

import (
	"crypto/hmac"
	"crypto/sha256"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lntypes"
)

// Oven bakes macaroons by combining the root secret, user ID, and caveats.
type Oven struct {
	uid     secrets.UserId
	root    secrets.Secret
	caveats []Caveat
	mac     *Macaroon
}

// NewOven creates a new Oven with the given root secret.
func NewOven(root secrets.Secret) Oven {
	oven := Oven{}
	oven.root = root
	return oven
}

// UserId sets the user ID in the Oven.
func (oven Oven) WithUserId(uid secrets.UserId) Oven {
	oven.uid = uid
	return oven
}

// Attenuate adds a single caveat to the Oven.
func (oven *Oven) Attenuate(caveat Caveat) {
	oven.caveats = append(oven.caveats, caveat)
}

// Adds multiple caveats to the Oven.
func (oven Oven) WithCaveats(caveats ...Caveat) Oven {
	oven.caveats = append(oven.caveats, caveats...)
	return oven
}

// Adds caveats for the specified services to the Oven.
func (oven Oven) WithService(services ...Service) Oven {
	for _, service := range services {
		oven.caveats = append([]Caveat{service.Caveat()}, oven.caveats...)
	}
	return oven
}

// Cook creates a new Macaroon by combining the user ID, caveats, and a signature.
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

	if oven.mac != nil {
		caveats = append(oven.mac.caveats, oven.caveats...)
	} else {
		caveats = oven.caveats
	}

	// Create and return the Macaroon with the user ID, caveats, and signature
	return Macaroon{user_id: oven.uid, caveats: caveats, signature: signature}, nil
}
