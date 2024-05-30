package main

import (
	"context"
	"fmt"
	"io"
	"lsat/challenge"
	"lsat/macaroon"
	"lsat/mock"
	"net/http"
	"strings"

	"github.com/lightningnetwork/lnd/lntypes"
)

const (
	authURL = "http://localhost:8080"
)

type TestClient struct{}

func main() {
	client := TestClient{}
	client.sendTokenRequest()
}

func parseToken(mac string, invoice string) (macaroon.Token, error) {
	Macaroon, err := decodeMacaroon(mac)
	if err != nil {
		return macaroon.Token{}, fmt.Errorf("error decoding macaroon: %v", err)
	}
	Preimage, err := decodePreimage(invoice)
	if err != nil {
		return macaroon.Token{}, fmt.Errorf("error decoding preimage: %v", err)
	}

	return macaroon.Token{Macaroon, Preimage}, nil
}

func decodeMacaroon(input string) (macaroon.Macaroon, error) {
	parts := strings.Split(input, "\"")
	if len(parts) < 2 {
		return macaroon.Macaroon{}, fmt.Errorf("invalid macaroon string")
	}

	macaroonStr := parts[1]
	return macaroon.DecodeBase64(macaroonStr)
}

func decodePreimage(input string) (lntypes.Preimage, error) {
	parts := strings.Split(input, "\"")
	if len(parts) < 2 {
		return lntypes.Preimage{}, fmt.Errorf("invalid preimage string")
	}

	invoiceStr := parts[1]
	ln := mock.TestLightningNode{Balance: 100000}
	payment, err := ln.PayInvoice(context.Background(), challenge.PayInvoiceRequest{Invoice: invoiceStr})
	if err != nil {
		return lntypes.Preimage{}, err
	}

	return payment.Preimage, nil
}

func (c *TestClient) sendTokenRequest() {
	fmt.Println("Requesting Token...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "L402")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusPaymentRequired {
		wwwAuthenticateHeader := resp.Header.Get("WWW-Authenticate")
		parts := strings.Split(wwwAuthenticateHeader, " ")

		if len(parts) > 2 && parts[0] == "L402" {
			macaroon := parts[1]
			invoice := parts[2]

			token, err := parseToken(macaroon, invoice)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(token.Macaroon.ToJSON())
			c.sendAuthorizationRequest(authURL, token)
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("(%s) Unexpected response status: %s\n", resp.Status, string(body))
	}
}

func (c *TestClient) sendAuthorizationRequest(address string, token macaroon.Token) {
	fmt.Println("Sending Authorization Request...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", address+"/protected", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("L402 %s", token.String()))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	} else {
		fmt.Println("Unexpected response status:", resp.Status)
	}
}
