package macaroon

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lntypes"
)

// Version is an alias for the Macaroon version.
type Version = int8

// Macaroon struct represents an LSAT (Lightning Service Authentication Token) macaroon.
type Macaroon struct {
	user_id   secrets.UserId
	caveats   []Caveat
	signature lntypes.Hash
}

// Uid returns the user ID associated with the macaroon.
func (mac *Macaroon) UserId() secrets.UserId {
	return mac.user_id
}

// Services extracts service names from the Macaroon's caveats.
func (mac *Macaroon) Services() ServiceIterator {
	return ServiceIterator{caveats: mac.caveats}
}

// Caveats returns the list of caveats associated with the macaroon.
func (mac *Macaroon) Caveats() []Caveat {
	return mac.caveats
}

// Signature returns the signature of the macaroon.
func (mac *Macaroon) Signature() string {
	return mac.signature.String()
}

// String returns the string representation of the macaroon.
func (mac Macaroon) String() string {
	// Marshal the Macaroon struct to JSON
	jsonData, _ := json.Marshal(mac.ToJSON())

	// Encode the JSON data to base64
	base64String := base64.StdEncoding.EncodeToString(jsonData)

	return base64String
}

// Create an oven from a Macaroon.
//
// This is used when adding third party caveats.
func (mac *Macaroon) Oven() Oven {
	root, _ := secrets.MakeSecret(mac.signature[:])
	return Oven{
		root: root,
		uid:  mac.user_id,
		mac:  mac,
	}
}

// macaroonJSON struct is used for JSON encoding/decoding of macaroon.
type macaroonJSON struct {
	Uid     string   `json:"user_id"`
	Caveats []Caveat `json:"caveats"`
	Sig     string   `json:"signature"`
}

// ToJSON converts Macaroon to macaroonJSON.
func (mac *Macaroon) ToJSON() macaroonJSON {
	return macaroonJSON{
		Uid:     mac.user_id.String(),
		Caveats: mac.caveats,
		Sig:     mac.Signature(),
	}
}

// decodeBase64 decodes a base64-encoded string.
func decodeBase64(encodedString string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

// Decode decodes a base64-encoded macaroon string into a Macaroon struct.
func DecodeBase64(encodedString string) (Macaroon, error) {
	// Decode the base64 string
	decoded, err := decodeBase64(encodedString)
	if err != nil {
		return Macaroon{}, err
	}

	// Unmarshal the decoded data into the macaroonJSON type
	var macJSON macaroonJSON
	err = json.Unmarshal(decoded, &macJSON)

	if err != nil {
		return Macaroon{}, err
	}

	// Convert the hex-encoded UID and signature to their respective types
	uid, _ := hex.DecodeString(macJSON.Uid)
	sig, _ := hex.DecodeString(macJSON.Sig)

	// Create lntypes.Hash and secrets.UserId from the decoded values
	sigHash, _ := lntypes.MakeHash(sig)
	uidHash, _ := secrets.MakeUserId(uid)

	// Create a Macaroon struct
	mac := Macaroon{
		user_id:   uidHash,
		caveats:   macJSON.Caveats,
		signature: sigHash,
	}

	return mac, nil
}
