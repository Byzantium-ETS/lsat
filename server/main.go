package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"lsat/secrets"

	"github.com/lightningnetwork/lnd/lntypes"
)

type Handler struct {
	*auth.Minter
}

const (
	macaroonHeader    = "L402"
	authFailedMessage = "Authentication failed!"
	redirectURL       = "https://picsum.photos/500"
)

var serviceName = getEnv("SERVICE_NAME", "image")

var (
	config = auth.NewConfig([]macaroon.Service{
		macaroon.NewService(serviceName, 1000),
	})
	secretStore = secrets.NewSecretFactory()
	challenger  = mock.NewChallenger()
)

func main() {
	host := getEnv("HOST", "localhost:8080")

	minter := auth.NewMinter(config, secretStore, challenger)
	handler := &Handler{Minter: &minter}

	log.Printf("Server launched at %s\n", host)
	http.HandleFunc("/", handler.handleAuthorization)
	http.HandleFunc("/preflight", handlePreflight)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

func handlePreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, WWW-Authenticate")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleAuthorization(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		handlePreflight(w, r)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != macaroonHeader {
		uid := secrets.NewUserId()

		pretoken, err := h.Minter.MintToken(uid, macaroon.NewServiceId(getEnv("SERVICE_NAME", "image"), 0))
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

	log.Printf("Authorization: %s", mac.UserId())

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
