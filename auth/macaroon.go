package auth

import "lsat/secrets"

type Macaroon struct {
	id      MacaroonId
	caveats []Caveat
	sig     string
}

type MacaroonId struct {
	version int8
	hash    int64
	uid     int32
}

type Oven struct {
	caveats []Caveat
}

func (oven Oven) Raw(root secrets.RootKey) (Macaroon, error) {
	return Macaroon{}, nil
}
