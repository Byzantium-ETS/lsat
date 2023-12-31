package macaroon

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lntypes"
)

type Version = int8

type Macaroon struct {
	uid     secrets.UserId
	caveats []Caveat
	sig     lntypes.Hash
}

type macaroonJSON struct {
	Uid     string   `json:"user_id"`
	Caveats []Caveat `json:"caveats"`
	Sig     string   `json:"signature"`
}

func (mac Macaroon) Uid() secrets.UserId {
	return mac.uid
}

func (mac Macaroon) Caveats() []Caveat {
	return mac.caveats
}

func (mac Macaroon) Signature() string {
	return mac.sig.String()
}

func (mac Macaroon) toJSON() macaroonJSON {
	return macaroonJSON{
		Uid:     mac.uid.String(),
		Caveats: mac.caveats,
		Sig:     mac.Signature(),
	}
}

func (mac Macaroon) String() string {
	// Marshal the Macaroon struct to JSON
	jsonData, err := json.Marshal(mac.toJSON())

	if err != nil {
		return fmt.Sprint(err)
	}

	// Encode the JSON data to base64
	base64String := base64.StdEncoding.EncodeToString(jsonData)

	return base64String
}

func decodeBase64(encodedString string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func Decode(encodedString string) (Macaroon, error) {
	// Decode the base64 string
	decoded, err := decodeBase64(encodedString)
	if err != nil {
		return Macaroon{}, err
	}

	// Unmarshal the decoded data into the Macaroon type
	var macJSON macaroonJSON
	err = json.Unmarshal(decoded, &macJSON)

	if err != nil {
		fmt.Println(macJSON.Caveats)
		return Macaroon{}, err
	}

	uid, _ := hex.DecodeString(macJSON.Sig)
	sig, _ := hex.DecodeString(macJSON.Sig)

	sig_hash, _ := lntypes.MakeHash(sig)
	uid_hash, _ := secrets.MakeUserId(uid)

	mac := Macaroon{
		uid:     uid_hash,
		caveats: macJSON.Caveats,
		sig:     sig_hash,
	}

	return mac, nil
}

// Bakes macaroons
type Oven struct {
	uid     secrets.UserId
	root    secrets.Secret
	caveats []Caveat
}

func NewOven(root secrets.Secret) Oven {
	oven := Oven{}
	oven.root = root
	return oven
}

func (oven Oven) UserId(uid secrets.UserId) Oven {
	oven.uid = uid
	return oven
}

func (oven Oven) Attenuate(caveat Caveat) Oven {
	oven.caveats = append(oven.caveats, caveat)
	return oven
}

func (oven Oven) Caveats(caveats ...Caveat) Oven {
	oven.caveats = append(oven.caveats, caveats...)
	return oven
}

func (oven Oven) Service(services ...Service) Oven {
	for _, service := range services {
		oven.caveats = append([]Caveat{service.Caveat()}, oven.caveats...)
	}
	return oven
}

func (oven Oven) Cook() (Macaroon, error) {
	// Je crois que c'est ca l'idee
	mac := hmac.New(sha256.New, oven.root[:])

	for _, caveat := range oven.caveats {
		mac.Write([]byte(caveat.String()))
	}

	signature, err := lntypes.MakeHash(mac.Sum(nil))

	if err != nil {
		return Macaroon{}, err
	}

	return Macaroon{uid: oven.uid, caveats: oven.caveats, sig: signature}, nil
}
