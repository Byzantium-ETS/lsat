package macaroon

import (
	"context"
	"lsat/lightning"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
)

const (
	BaseVersion = iota
)

// A service token.
//
// It holds the macaroon and its secret.
type Token struct {
	Mac      Macaroon
	Preimage lntypes.Preimage // Le secret de la transaction
}

// A transitive service token.
//
// It needs to be paid in order to become effective.
type PreToken struct {
	Mac     Macaroon
	Invoice lnrpc.AddInvoiceResponse // L'invoice qui sera payé par le client
}

// Create a Token.
func (token PreToken) Pay(node lightning.InvoiceHandler) (Token, error) {
	preimage, err := node.Pay(context.Background(), token.Invoice)
	if err != nil {
		return Token{Mac: token.Mac, Preimage: preimage}, nil
	} else {
		return Token{}, err
	}
}

// A key used to identify macaroons in the database.
type TokenID struct {
	version Version
	uid     secrets.UserId
	hash    lntypes.Hash // Le hash du preimage de la transaction
}

func NewTokenID(uid secrets.UserId, hash lntypes.Hash) TokenID {
	return TokenID{version: BaseVersion, uid: uid, hash: hash}
}
