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
	host = "localhost:8080"

	macaroonHeader    = "L402"
	serviceName       = "image"
	authFailedMessage = "Authentication failed!"
)

var (
	serviceLimiter = auth.NewConfig([]macaroon.Service{
		macaroon.NewService(serviceName, 1000),
	})
	secretStore = secrets.NewSecretFactory()
	challenger  = mock.NewChallenger()
)

func main() {
	// Initialize your Server instance
	minter := auth.NewMinter(serviceLimiter, &secretStore, challenger)

	// Create a Handler with access to the Minter
	handle := &Handler{Minter: &minter}

	// Launch the server
	fmt.Println("Server launched at", host)
	http.HandleFunc("/", handle.handleAuthorization)
	http.HandleFunc("/protected", handle.handleProtected)
	err := http.ListenAndServe(host, nil)
	fmt.Println(err)
}

func (h *Handler) handleAuthorization(w http.ResponseWriter, r *http.Request) {
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
		pretoken, err := h.Minter.MintToken(uid, macaroon.NewServiceId(serviceName, 0))

		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "%s", err)
			return
		}

		macaroon := pretoken.Macaroon

		// Format Macaroon and invoice in WWW-Authenticate header
		authHeader := fmt.Sprintf("%s macaroon=\"%s\", invoice=\"%s\"", macaroonHeader, macaroon, pretoken.InvoiceResponse.Invoice)

		// Set the WWW-Authenticate header
		w.Header().Set("WWW-Authenticate", authHeader)

		w.WriteHeader(http.StatusPaymentRequired)
		fmt.Fprint(w, "Payment Required")
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "%s", authFailedMessage)
}

func (h *Handler) handleProtected(w http.ResponseWriter, r *http.Request) {
	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Check if Authorization header is present
	// Parse the Authorization header
	parts := strings.Split(authHeader, " ")

	credentials := strings.Split(parts[1], ":")

	Macaroon, _ := macaroon.DecodeBase64(credentials[0])
	Preimage, _ := lntypes.MakePreimageFromStr(credentials[1])

	token := macaroon.Token{
		Macaroon: Macaroon,
		Preimage: Preimage,
	}

	err := h.Minter.AuthToken(&token)

	if err == nil {
		// Return an image
		http.Redirect(w, r, "https://picsum.photos/500", http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%s", err)
	}
}
