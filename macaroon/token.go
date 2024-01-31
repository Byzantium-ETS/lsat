package macaroon

import (
	"context"
	"encoding/hex"
	"fmt"
	"lsat/challenge"
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
	Macaroon Macaroon         // The macaroon.
	Preimage lntypes.Preimage // The secret of the transaction.
}

func (token Token) String() string {
	// Encode the Macaroon(s) as base64
	macaroonBase64 := token.Macaroon.String()

	// Encode the Preimage as hex
	preimageHex := hex.EncodeToString(token.Preimage[:])

	// Combine the encoded Macaroon(s) and encoded Preimage as <macaroon(s)>:<preimage>
	encodedToken := fmt.Sprintf("%s:%s", macaroonBase64, preimageHex)

	return encodedToken
}

// A transitive service token.
//
// It needs to be paid in order to become effective.
//
// This object is sent when the Macaroon is minted.
type PreToken struct {
	Macaroon       Macaroon                 // The macaroon.
	PaymentRequest lnrpc.AddInvoiceResponse // The invoice sent to the user.
}

// Pay a token.
//
// This creates a valid Token.
func (token PreToken) Pay(node challenge.LightningNode) (Token, error) {
	cx := context.Background()
	cx = context.WithValue(cx, "macaroon", token.Macaroon) // Enrich the context with a macaroon
	preimage, err := node.SendPayment(cx, token.PaymentRequest)
	if err != nil {
		return Token{Macaroon: token.Macaroon, Preimage: preimage}, nil
	} else {
		return Token{}, err
	}
}

func (token PreToken) String() string {
	// Encode the Macaroon(s) as base64
	macaroonBase64 := token.Macaroon.String()

	// Encode the Invoice
	invoice := token.PaymentRequest.GetPaymentRequest()

	// Combine the encoded Macaroon(s) and encoded Preimage as <macaroon(s)>:<preimage>
	encodedToken := fmt.Sprintf("%s:%s", macaroonBase64, invoice)

	return encodedToken
}

// A key used to identify macaroons in the database.
type TokenID struct {
	version Version
	user_id secrets.UserId // The id of the token owner
	hash    lntypes.Hash   // The hash of the preimage of the transaction
}

func (token Token) Id() TokenID {
	return TokenID{
		version: BaseVersion,
		user_id: token.Macaroon.UserId(),
		hash:    lntypes.Hash(token.Preimage),
	}
}
