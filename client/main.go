package main

import (
	"fmt"
	"io"
	"lsat/macaroon"
	"net/http"
	"strings"

	"github.com/lightningnetwork/lnd/lntypes"
)

const (
	serverAddress = "http://localhost:8080"
)

type TestClient struct {
	macaroon macaroon.Macaroon
}

func main() {
	client := TestClient{}
	client.sendAuthenticationRequest()
}

func parseLSATString(mac string, invoice string) (macaroon.Token, error) {
	Macaroon, err := parseMacaroonString(mac)
	if err != nil {
		return macaroon.Token{}, fmt.Errorf("Error decoding macaroon: %v", err)
	}
	Preimage, err := parsePreimageString(invoice)
	if err != nil {
		return macaroon.Token{}, fmt.Errorf("Error decoding preimage: %v", err)
	}

	// You can use the decoded macaroon and invoice as needed in your application
	// For example, you might want to return the macaroon
	return macaroon.Token{
		Macaroon: Macaroon,
		Preimage: Preimage,
	}, nil
}

func parseMacaroonString(input string) (macaroon.Macaroon, error) {
	// Split the input string into parts
	parts := strings.Split(input, "\"")

	macaroonStr := parts[1]

	return macaroon.DecodeBase64(macaroonStr)
}

func parsePreimageString(input string) (lntypes.Preimage, error) {
	// Split the input string into parts
	parts := strings.Split(input, "\"")

	invoiceStr := parts[1]

	return lntypes.MakePreimageFromStr(invoiceStr)
}

// TO-DO Put the response in a struct
func (c *TestClient) sendAuthenticationRequest() {
	fmt.Println("Sending Authentication Request...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", serverAddress, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the LSAT Authorization header with an invalid token for demonstration
	req.Header.Set("Authorization", "L402")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the response has a 402 status code
	if resp.StatusCode == http.StatusPaymentRequired {
		// Parse the WWW-Authenticate header to get the LSAT token and invoice
		wwwAuthenticateHeader := resp.Header.Get("WWW-Authenticate")
		parts := strings.Split(wwwAuthenticateHeader, " ")

		if len(parts) > 1 && parts[0] == "L402" {
			macaroon := parts[1]
			invoice := parts[2]

			token, err := parseLSATString(macaroon, invoice)

			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(token.Macaroon.ToJSON(), "\n")

			c.sendAuthorizationRequest(token)
		}
	} else {
		fmt.Println("Unexpected response status:", resp.Status)
	}
}

func (c *TestClient) sendAuthorizationRequest(token macaroon.Token) {
	fmt.Println("Sending Authorization Request...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", serverAddress, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the LSAT Authorization header with an invalid token for demonstration
	req.Header.Set("Authorization", fmt.Sprintf("L402 %s", token.String()))

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the response has a 402 status code
	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		// Parse the WWW-Authenticate header to get the LSAT token and invoice
		encodedMac := string(body)

		if len(encodedMac) > 0 {
			signedMac, err := macaroon.DecodeBase64(encodedMac)
			if err != nil {
				return
			}

			fmt.Println(signedMac.ToJSON(), "\n")

			c.macaroon = signedMac

			c.sendProtectedRequest()
		}
	} else {
		fmt.Println("Unexpected response status:", resp.Status)
	}
}

// Reuse the response from the first GET
func (c *TestClient) sendProtectedRequest() {
	fmt.Println("Sending Protected Request...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", serverAddress+"/protected", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the LSAT Authorization header with a valid token and preimage for demonstration
	req.Header.Set("Authorization", fmt.Sprintf("L402 %s", c.macaroon))

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Request unauthorized!")
		return
	}

	defer resp.Body.Close()

	fmt.Println(resp.Status)
}
