package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lntypes"
)

type Version = int8

type Macaroon struct {
	caveats []Caveat
	sig     string
}

func (mac Macaroon) Caveats() []Caveat {
	return mac.caveats
}

func (mac Macaroon) Signature() string {
	return mac.Signature()
}

// La clé utlisée pour map les macaroons dans la base de données.
type MacaroonId struct {
	version Version
	hash    lntypes.Hash
	Uid     secrets.UserId
}

// Bakes macaroons
type Oven struct {
	root    secrets.Secret
	caveats []Caveat
}

func NewOven(store *secrets.SecretStore, uid secrets.UserId) (Oven, error) {
	return Oven{}, nil
}

func (oven Oven) Attenuate(caveat Caveat) Oven {
	oven.caveats = append(oven.caveats, caveat)
	return oven
}

func (oven Oven) MapCaveats(caveats []Caveat) Oven {
	oven.caveats = append(oven.caveats, caveats...)
	return oven
}

func (oven Oven) Cook() (Macaroon, error) {
	// Je crois que c'est ca l'idee
	mac := hmac.New(sha256.New, oven.root)
	for _, v := range oven.caveats {
		mac = hmac.New(func() hash.Hash { return mac }, []byte(v.ToString()))
	}
	return Macaroon{caveats: oven.caveats, sig: string(mac.Sum(nil))}, nil
}
