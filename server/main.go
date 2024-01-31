package main

import (
	"fmt"
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"lsat/secrets"
	"net/http"
	"strings"

	"github.com/lightningnetwork/lnd/lntypes"
)

type Handler struct {
	*auth.Minter
}

const (
	address           = "localhost:8080"
	macaroonHeader    = "L402"
	catService        = mock.CatService
	authFailedMessage = "Authentication failed!"
)

var (
	serviceLimiter = mock.NewServiceLimiter()
	secretStore    = mock.NewTestStore()
	challenger     = mock.NewChallenger()
)

func main() {
	// Initialize your Server instance
	minter := auth.NewMinter(&serviceLimiter, &secretStore, challenger)

	// Create a Handler with access to the Minter
	handle := &Handler{Minter: &minter}

	fmt.Println("Server launched at", address)
	http.HandleFunc("/", handle.handleAuthentication)
	http.HandleFunc("/protected", handle.handleAuthorization)
	err := http.ListenAndServe(address, nil)
	fmt.Println(err)
}

func (h *Handler) handleAuthorization(w http.ResponseWriter, r *http.Request) {
	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Check if Authorization header is present
	// Parse the Authorization header
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != macaroonHeader {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unknown request!")
		return
	}

	macaroon, _ := macaroon.DecodeBase64(parts[1])

	err := serviceLimiter.VerifyMacaroon(&macaroon)

	if err == nil {
		// Respond with success (for demonstration purposes)
		// We should respond with the resource
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Request authorized!")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%s", authFailedMessage)
	}
}

func (h *Handler) handleAuthentication(w http.ResponseWriter, r *http.Request) {
	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Check if Authorization header is present
	// Parse the Authorization header
	parts := strings.Split(authHeader, " ")

	// Extract the Macaroon and preimage from the Authorization header
	if len(parts) != 2 || parts[0] != macaroonHeader {
		// Create a new UserId
		uid := secrets.NewUserId()

		// Invalid Authorization header format, respond with 402 Payment Required and WWW-Authenticate header
		pretoken, err := h.Minter.MintToken(uid, catService)

		if err != nil {
			fmt.Println(err)
			return
		}

		macaroon := pretoken.Macaroon

		// Format Macaroon and invoice in WWW-Authenticate header
		authHeader := fmt.Sprintf("%s macaroon=\"%s\", invoice=\"%s\"", macaroonHeader, macaroon, pretoken.PaymentRequest.GetPaymentRequest())

		// Set the WWW-Authenticate header
		w.Header().Set("WWW-Authenticate", authHeader)

		w.WriteHeader(http.StatusPaymentRequired)
		fmt.Fprint(w, "Payment Required")
		return
	}

	credentials := strings.Split(parts[1], ":")

	Macaroon, _ := macaroon.DecodeBase64(credentials[0])
	Preimage, _ := lntypes.MakePreimageFromStr(credentials[1])

	token := macaroon.Token{
		Macaroon: Macaroon,
		Preimage: Preimage,
	}

	// Process the request with the extracted Macaroon and preimage
	// Your logic for handling the request goes here...
	signedMac, err := h.Minter.AuthToken(&token)

	if err == nil {
		// Respond with success (for demonstration purposes)
		// We should respond with the resource
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", signedMac)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%s %s", authFailedMessage, err)
	}
}
