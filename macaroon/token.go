package macaroon

import (
	"lsat/lightning"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lntypes"
)

const (
	BaseVersion = iota
)

// Un Token complet
type Token struct {
	Mac      Macaroon
	Preimage lntypes.Preimage // Le secret de la transaction
}

// Un Token partiel.
// Le Token complet est créé quand l'invoice est payé quand l'utilisateur
type PreToken struct {
	Mac     Macaroon
	Invoice string // L'invoice qui sera payé par le client
}

// Créé un Token.
// Utilisé par le client
func (token PreToken) Pay(node lightning.Node) (Token, error) {
	preimage, err := node.Pay(token.Invoice)
	if err != nil {
		return Token{Mac: token.Mac, Preimage: preimage}, nil
	} else {
		return Token{}, err
	}
}

// La clé utlisée pour map les macaroons dans la base de données.
type TokenID struct {
	version Version
	uid     secrets.UserId
	hash    lntypes.Hash // Le hash du preimage de la transaction
}

func NewTokenID(uid secrets.UserId, hash lntypes.Hash) TokenID {
	return TokenID{version: BaseVersion, uid: uid, hash: hash}
}
