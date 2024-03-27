package main

import (
	"fmt"
	"io"
	"log"
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"lsat/secrets"
	"net/http"
	"strings"
)

const (
	address           = "localhost:8443"
	authFailedMessage = "Authentication failed!"
)

var serviceLimiter = mock.NewServiceLimiter()

type Handler struct {
	mock.TestStore
}

func main() {
	handler := Handler{}
	fmt.Println("Server launched at", address)
	http.HandleFunc("/", handler.handleSecret)
	http.HandleFunc("/protected", handler.handleProtected)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (h *Handler) handleSecret(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	secret, err := secrets.MakeSecret(body)
	if err != nil {
		log.Fatal(err)
	}

	h.TestStore = mock.NewTestStoreFromSecret(secret)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleProtected(w http.ResponseWriter, r *http.Request) {
	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Check if Authorization header is present
	// Parse the Authorization header
	parts := strings.Split(authHeader, " ")

	credentials := strings.Split(parts[1], ":")

	Macaroon, _ := macaroon.DecodeBase64(credentials[0])

	minter := auth.NewMinter(serviceLimiter, &h.TestStore, nil)

	// Process the request with the extracted Macaroon and preimage
	// Your logic for handling the request goes here...
	err := minter.AuthMacaroon(&Macaroon)

	if err == nil {
		// Respond with success (for demonstration purposes)
		// We should respond with the resource
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", "hello world!")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%s %s", authFailedMessage, err)
	}
}
