package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
	"lsat/secrets"
)

type Version = int8

type Macaroon struct {
	caveats []Caveat
	sig     string
}

func (mac Macaroon) Caveats() []Caveat {
	return mac.caveats
}

// La clé utlisée pour map les macaroons dans la base de données.
type MacaroonId struct {
	version Version
	hash    int64
	Uid     int32
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

func (oven Oven) Cook() (Macaroon, error) {
	// Je crois que c'est ca l'idee
	mac := hmac.New(sha256.New, oven.root)
	for _, v := range oven.caveats {
		mac = hmac.New(func() hash.Hash { return mac }, []byte(v.ToString()))
	}
	return Macaroon{caveats: oven.caveats, sig: string(mac.Sum(nil))}, nil
}
