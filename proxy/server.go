package proxy

import (
	"fmt"
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"lsat/secrets"
	"net/http"
	"strings"
)

const (
	service_name = "dogs"
)

var serviceManager = mock.TestServiceManager{}
var secretStore = mock.NewTestStore()
var challenger = mock.TestChallenger{}

type Handler struct {
	*auth.Minter
}

func HttpServer() {
	// Initialize your Server instance
	minter := auth.NewMinter(&serviceManager, &secretStore, &challenger)

	// Create a Handler with access to the Server
	handle := &Handler{Minter: &minter}

	fmt.Println("Server launched!")
	http.HandleFunc("/", handle.handleRequest) // Retourné dans le cas où la platforme reçoit un Token invalide
	http.ListenAndServe(":8080", nil)          // Ca devrait etre lié à la platforme?
}

func (h *Handler) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Create a new UserId
	uid := secrets.NewUserId()

	// Check if Authorization header is present
	// Parse the Authorization header
	parts := strings.Split(authHeader, " ")
	// Extract the Macaroon and preimage from the Authorization header
	credentials := strings.Split(parts[1], ":")
	if authHeader == "" || len(parts) != 2 || parts[0] != "L402" || len(credentials) != 2 {
		// Invalid Authorization header format, respond with 402 Payment Required and WWW-Authenticate header
		pretoken, err := h.Minter.MintToken(uid, service_name)

		if err != nil {
			fmt.Println(err)
			return
		}

		macaroon := pretoken.Mac.String()

		// Format Macaroon and invoice in WWW-Authenticate header
		authHeader := fmt.Sprintf("L402 macaroon=\"%s\", invoice=\"%s\"", macaroon, &pretoken.Invoice)

		// Set the WWW-Authenticate header
		w.Header().Set("WWW-Authenticate", authHeader)

		w.WriteHeader(http.StatusPaymentRequired)
		fmt.Fprint(w, "Payment Required")
		return
	}

	preimage := credentials[1]

	encodeding := credentials[0] // Encode Macaroon
	macaroon, err := macaroon.Decode(encodeding)

	if err != nil {
		fmt.Println(err)
		return
	}

	// token = macaroon.Token{}

	// Process the request with the extracted Macaroon and preimage
	// Your logic for handling the request goes here...
	// err = h.Minter.VerifyMacaroon(&macaroon)

	if err == nil {
		// Respond with success (for demonstration purposes)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Request authorized with Macaroon: %s and Preimage: %s", macaroon, preimage)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Authentification failed! %s", err)
	}

}
