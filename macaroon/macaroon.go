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
	userId    secrets.UserID
	caveats   []Caveat
	signature lntypes.Hash
}

// Uid returns the user ID associated with the macaroon.
func (mac *Macaroon) UserId() secrets.UserID {
	return mac.userId
}

// func (mac *Macaroon) Services() service ServiceIterator {
// 	return ServiceIterator{caveats: mac.caveats}
// }

// Caveats returns the list of caveats associated with the macaroon.
func (mac *Macaroon) Caveats() []Caveat {
	return mac.caveats
}

// Signature returns the signature of the macaroon.
func (mac *Macaroon) Signature() lntypes.Hash {
	return mac.signature
}

// Returns the Value of the caveat with the given Key
func (mac *Macaroon) GetValue(key string) ValueIterator {
	return NewIterator(key, mac.caveats)
}

func (mac Macaroon) String() string {
	// Marshal the Macaroon struct to JSON
	jsonData, err := json.Marshal(mac.ToJSON())

	if err != nil {
		panic(err)
	}

	// Encode the JSON data to base64
	base64String := base64.StdEncoding.EncodeToString(jsonData)

	return base64String
}

// Create an oven from a Macaroon.
//
// This is used for adding third party caveats.
func (mac *Macaroon) Oven() Oven {
	root, _ := secrets.MakeSecret(mac.signature[:])
	return Oven{
		root:     root,
		userId:   mac.userId,
		macaroon: mac,
	}
}

// MacaroonJSON struct is used for JSON encoding/decoding of macaroon.
type MacaroonJSON struct {
	UserId    string   `json:"user_id"`
	Caveats   []Caveat `json:"caveats"`
	Signature string   `json:"signature"`
}

// ToJSON converts Macaroon to macaroonJSON.
func (mac *Macaroon) ToJSON() MacaroonJSON {
	return MacaroonJSON{
		UserId:    mac.userId.String(),
		Caveats:   mac.caveats,
		Signature: mac.Signature().String(),
	}
}

func (mac MacaroonJSON) String() string {
	// Marshal the Macaroon struct to JSON
	jsonData, err := json.MarshalIndent(mac, "", "  ")

	if err != nil {
		panic(err)
	}

	return string(jsonData)
}

// Unwrap get a Macaroon from the JSON object.
func (mac MacaroonJSON) Unwrap() (Macaroon, error) {
	signature, err := lntypes.MakeHashFromStr(mac.Signature)
	if err != nil {
		return Macaroon{}, err
	}

	userId, err := secrets.MakeUserIdFromStr(mac.UserId)
	if err != nil {
		return Macaroon{}, err
	}

	return Macaroon{
		userId:    userId,
		signature: signature,
		caveats:   mac.Caveats,
	}, nil
}

func (mac *Macaroon) MarshalJSON() ([]byte, error) {
	return json.Marshal(mac.ToJSON())
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
	var macJSON MacaroonJSON
	err = json.Unmarshal(decoded, &macJSON)

	if err != nil {
		return Macaroon{}, err
	}

	// Convert the hex-encoded UID and signature to their respective types
	uid, _ := hex.DecodeString(macJSON.UserId)
	sig, _ := hex.DecodeString(macJSON.Signature)

	// Create lntypes.Hash and secrets.UserId from the decoded values
	sigHash, _ := lntypes.MakeHash(sig)
	uidHash, _ := secrets.MakeUserId(uid)

	// Create a Macaroon struct
	mac := Macaroon{
		userId:    uidHash,
		caveats:   macJSON.Caveats,
		signature: sigHash,
	}

	return mac, nil
}
