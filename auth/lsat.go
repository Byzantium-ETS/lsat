package auth

import (
	"lsat/lightning"

	"github.com/lightningnetwork/lnd/lntypes"
)

// Un LSAT complet
type Token struct {
	Mac       Macaroon
	Service   string           // Le nom du service activé par le token
	Pre_image lntypes.Preimage // Le secret de la transaction
}

// Un LSAT partiel.
// Le LSAT complet est créé quand l'invoice est payé quand l'utilisateur
type preToken struct {
	Mac     Macaroon
	Service string // Le nom du service activé par le token
	Invoice string // L'invoice qui sera payé par le client
}

// Créé un Token.
// Utilisé par le client
func (token preToken) Pay(node lightning.LightningNode) (Token, error) {
	preimage, err := node.Pay(token.Invoice)
	if err != nil {
		return Token{Mac: token.Mac, Service: token.Service, Pre_image: preimage}, nil
	} else {
		return Token{}, err
	}
}
