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
	Mac      Macaroon         // The macaroon.
	Preimage lntypes.Preimage // The secret of the transaction.
}

func (token Token) String() string {
	// Encode the Macaroon(s) as base64
	macaroonBase64 := token.Mac.String()

	// Encode the Preimage as hex
	preimageHex := hex.EncodeToString(token.Preimage[:])

	// Combine the encoded Macaroon(s) and encoded Preimage as <macaroon(s)>:<preimage>
	encodedToken := fmt.Sprintf("%s:%s", macaroonBase64, preimageHex)

	return encodedToken
}

const (
	AuthSchemeErr = "The authentication scheme is not L402!"
)

// A transitive service token.
//
// It needs to be paid in order to become effective.
//
// This object is sent when the Macaroon is minted.
type PreToken struct {
	Mac     Macaroon                 // The macaroon.
	Invoice lnrpc.AddInvoiceResponse // The invoice sent to the user.
}

// Create a Token.
func (token PreToken) Pay(node challenge.InvoiceHandler) (Token, error) {
	cx := context.Background()
	cx = context.WithValue(cx, "macaroon", token.Mac) // Enrich the context with a macaroon
	preimage, err := node.SendPayment(cx, token.Invoice)
	if err != nil {
		return Token{Mac: token.Mac, Preimage: preimage}, nil
	} else {
		return Token{}, err
	}
}

func (token PreToken) String() (string, error) {
	// Encode the Macaroon(s) as base64
	macaroonBase64 := token.Mac.String()

	// Encode the Invoice
	invoice := token.Invoice.PaymentRequest

	// Combine the encoded Macaroon(s) and encoded Preimage as <macaroon(s)>:<preimage>
	encodedToken := fmt.Sprintf("%s:%s", macaroonBase64, invoice)

	return encodedToken, nil
}

// A key used to identify macaroons in the database.
type TokenID struct {
	version Version
	uid     secrets.UserId // The id of the token owner
	hash    lntypes.Hash   // The hash of the preimage of the transaction
}

func (token Token) Id() TokenID {
	return TokenID{
		version: BaseVersion,
		uid:     token.Mac.Uid(),
		hash:    lntypes.Hash(token.Preimage),
	}
}
