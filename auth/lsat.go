package auth

type Token struct {
	mac       Macaroon
	pre_image string
}
