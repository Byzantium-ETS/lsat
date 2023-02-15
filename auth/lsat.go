package auth

import (
	"lsat/lightning"

	"github.com/lightningnetwork/lnd/lntypes"
)

// Un LSAT complet
type LSAT struct {
	Mac      Macaroon
	Preimage lntypes.Preimage // Le secret de la transaction
}

// Un LSAT partiel.
// Le LSAT complet est créé quand l'invoice est payé quand l'utilisateur
type preLSAT struct {
	Mac     Macaroon
	Invoice string // L'invoice qui sera payé par le client
}

// Créé un LSAT.
// Utilisé par le client
func (token preLSAT) Pay(node lightning.LightningNode) (LSAT, error) {
	preimage, err := node.Pay(token.Invoice)
	if err != nil {
		return LSAT{Mac: token.Mac, Preimage: preimage}, nil
	} else {
		return LSAT{}, err
	}
}
