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
	host              = "localhost:8080"
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
	minter := auth.NewMinter(serviceLimiter, &secretStore, challenger)
	handler := &Handler{Minter: &minter}

	fmt.Println("Server launched at", host)
	http.HandleFunc("/", handler.handleAuthorization)
	http.HandleFunc("/protected", handler.handleProtected)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) handleAuthorization(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != macaroonHeader {
		uid := secrets.NewUserId()

		pretoken, err := h.Minter.MintToken(uid, macaroon.NewServiceId(serviceName, 0))
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "%s", err)
			return
		}

		mac := pretoken.Macaroon
		authHeader := fmt.Sprintf("%s macaroon=\"%s\", invoice=\"%s\"", macaroonHeader, mac, pretoken.InvoiceResponse.Invoice)

		w.Header().Set("WWW-Authenticate", authHeader)
		w.WriteHeader(http.StatusPaymentRequired)
		fmt.Fprint(w, "Payment Required")
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "%s", authFailedMessage)
}

func (h *Handler) handleProtected(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != macaroonHeader {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%s", authFailedMessage)
		return
	}

	credentials := strings.Split(parts[1], ":")
	if len(credentials) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%s", authFailedMessage)
		return
	}

	mac, err := macaroon.DecodeBase64(credentials[0])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%s", err)
		return
	}

	preimage, err := lntypes.MakePreimageFromStr(credentials[1])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%s", err)
		return
	}

	token := macaroon.Token{
		Macaroon: mac,
		Preimage: preimage,
	}

	err = h.Minter.AuthToken(&token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%s", err)
		return
	}

	http.Redirect(w, r, "https://picsum.photos/500", http.StatusOK)
}
