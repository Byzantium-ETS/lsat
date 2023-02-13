package auth

import "lsat/lightning"

type Token struct {
	mac       Macaroon
	pre_image string
}

type preToken struct {
	mac     Macaroon
	invoice string
}

func (token Token) Pay(node lightning.LightningNode) Token {
	return Token{}
}
