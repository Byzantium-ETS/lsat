package macaroon

import (
	"crypto/hmac"
	"crypto/sha256"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lntypes"
)

type Version = int8

type Macaroon struct {
	caveats []Caveat
	service Service
	sig     lntypes.Hash
}

func (mac Macaroon) Caveats() []Caveat {
	return mac.caveats
}

func (mac Macaroon) Signature() string {
	return mac.sig.String()
}

func (mac Macaroon) Service() Service {
	return mac.service
}

// Bakes macaroons
type Oven struct {
	service Service
	root    secrets.Secret
	caveats []Caveat
}

func NewOven(root secrets.Secret) Oven {
	oven := Oven{}
	oven.root = root
	return oven
}

func (oven Oven) Attenuate(caveat Caveat) Oven {
	oven.caveats = append(oven.caveats, caveat)
	return oven
}

func (oven Oven) MapCaveats(caveats []Caveat) Oven {
	oven.caveats = append(oven.caveats, caveats...)
	return oven
}

func (oven Oven) Service(service Service) Oven {
	oven.service = service
	return oven
}

func (oven Oven) Cook() (Macaroon, error) {
	// Je crois que c'est ca l'idee
	mac := hmac.New(sha256.New, oven.root[:])

	// fmt.Println(oven.root)

	services, err := FmtServices(oven.service)

	if err == nil {
		mac.Write([]byte(services))
		// mac = hmac.New(func() hash.Hash { return mac }, services)
	} else {
		return Macaroon{}, err
	}

	for _, v := range oven.caveats {
		mac.Write([]byte(v.String()))
	}

	signature, err := lntypes.MakeHash(mac.Sum(nil))

	if err != nil {
		return Macaroon{}, err
	}

	return Macaroon{caveats: oven.caveats, service: oven.service, sig: signature}, nil
}
