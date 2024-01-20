package mock

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	. "lsat/macaroon"
	"lsat/secrets"
	"time"

	"github.com/lightningnetwork/lnd/lntypes"
)

const (
	DogService = "dogs"
	CatService = "cats"

	timeKey      = "time"
	signatureKey = "signature"

	timeErr = "the macaroon is expired!"
	signErr = "the signature of that token is invalid!"
)

type TestServiceLimiter struct {
	secret secrets.Secret
}

func NewServiceLimiter() TestServiceLimiter {
	return TestServiceLimiter{
		secret: secrets.NewSecret(),
	}
}

// listCaveats generates a list of caveats based on the provided service.
func listCaveats(service Service) []Caveat {
	// Use a switch statement to handle different services and their respective caveats.
	switch service.Name {
	case DogService, CatService:
		// If the service is a DogService or CatService, include a time-based caveat.
		// The caveat represents an expiration time one hour from now.
		return []Caveat{NewCaveat(timeKey, time.Now().Add(time.Hour).Format(time.Layout))}
	default:
		// If the service is not explicitly handled, return an empty slice.
		return nil
	}
}

// verifyCaveats checks the validity of the provided caveats.
func (s *TestServiceLimiter) verifyCaveats(caveats ...Caveat) error {
	for _, caveat := range caveats {
		switch caveat.Key {
		case timeKey:
			// Parse the value of the time caveat as a time.Time.
			expiry, err := time.Parse(time.Layout, caveat.Value)

			// If there is an error parsing the time, return the error.
			if err != nil {
				return err
			}

			// Check if the expiry time is before the current time.
			if expiry.Before(time.Now()) {
				return errors.New(timeErr)
			}
		}
	}
	// If all checks pass, return nil (no error).
	return nil
}

func (s *TestServiceLimiter) Services(cx context.Context, names ...string) ([]Service, error) {
	list := make([]Service, 0, len(names))
	for _, name := range names {
		switch name {
		case CatService:
			list = append(list, NewService(CatService, 1000))
		case DogService:
			list = append(list, NewService(DogService, 2000))
		default:
			return []Service{}, errors.New("unkown service!")
		}
	}
	return list, nil
}

func (s *TestServiceLimiter) Capabilities(cx context.Context, services ...Service) ([]Caveat, error) {
	arr := make([]Caveat, 0, len(services))
	for _, service := range services {
		arr = append(arr, listCaveats(service)...)
	}
	return arr, nil
}

// Sign signs the given macaroon by encrypting its signature with the service's secret.
// It adds a caveat containing the encrypted signature to the macaroon and returns the
// newly signed macaroon.
func (s *TestServiceLimiter) Sign(mac Macaroon) (Macaroon, error) {
	// Encrypt the macaroon's signature using the service's secret.
	signature, _ := encrypt(s.secret[:], mac.Signature().String())

	// Encode the encrypted signature in hexadecimal.
	encodedSig := hex.EncodeToString(signature)

	// Add a new caveat to the macaroon containing the encrypted signature.
	return mac.Oven().WithCaveats(NewCaveat(signatureKey, encodedSig)).Cook()
}

// VerifyMacaroon verifies the integrity and authenticity of the given macaroon.
// It checks the signature and validates the caveats.
func (s *TestServiceLimiter) VerifyMacaroon(mac *Macaroon) error {
	var signature lntypes.Hash
	var caveats []Caveat

	// Iterate through the macaroon's caveats to find the signature caveat.
	for i, caveat := range mac.Caveats() {
		if caveat.Key == signatureKey {
			decodedSig, err := hex.DecodeString(caveat.Value)
			if err != nil {
				return err
			}
			// Decrypt the signature value using the secret.
			rawSignature, err := decrypt(s.secret[:], decodedSig)
			if err != nil {
				return err
			}
			signature, _ = lntypes.MakeHashFromStr(rawSignature)
			caveats = mac.Caveats()[i:]
			break
		}
	}

	// Create a new secret using the hash as the root.
	root, err := secrets.MakeSecret(signature[:])
	if err != nil {
		return err
	}

	// Create a new macaroon with the extracted caveats and the new secret.
	newMac, err := NewOven(root).WithCaveats(caveats...).Cook()

	// Compare the signatures of the original and newly created macaroons.
	if newMac.Signature() != mac.Signature() {
		return errors.New(signErr)
	}

	// Verify the remaining caveats.
	return s.verifyCaveats(mac.Caveats()...)
}

func encrypt(secretKey []byte, plaintext string) ([]byte, error) {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return []byte{}, err
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return []byte{}, err
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	cipher := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return cipher, nil
}

func decrypt(secretKey []byte, ciphertext []byte) (string, error) {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, _ := gcm.Open(nil, []byte(nonce), ciphertext, nil)

	return string(plaintext), nil
}
