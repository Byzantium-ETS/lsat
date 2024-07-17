package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"lsat/auth"
	"lsat/challenge"
	"lsat/macaroon"
	"lsat/mock"
	"net/http"
	"strings"
	"time"
)

const (
	authURL = "http://localhost:8080"
)

type TestClient struct {
	tokenPath string
}

func main() {
	client := TestClient{
		tokenPath: getTokenPath(),
	}
	if client.tokenPath == "" {
		client.sendTokenRequest()
	} else {
		store, _ := auth.NewStore("./.store")
		token, err := store.GetTokenFromPath(client.tokenPath)
		if err != nil {
			log.Fatal(err)
		}
		client.sendAuthorizationRequest(authURL, *token)
	}
}

var lightningNode = mock.TestLightningNode{Balance: 10000}

// Connect to the phoenix node
// var lightningNode = phoenixd.NewPhoenixClient("baseUrl", "password")

func parsePreToken(mac string, invoice string) (macaroon.PreToken, error) {
	Macaroon, err := decodeMacaroon(mac)
	if err != nil {
		return macaroon.PreToken{}, fmt.Errorf("error decoding macaroon: %v", err)
	}
	return macaroon.PreToken{
		Macaroon: Macaroon,
		InvoiceResponse: challenge.InvoiceResponse{
			Invoice: strings.Split(invoice, "\"")[1],
		},
	}, nil
}

func decodeMacaroon(input string) (macaroon.Macaroon, error) {
	parts := strings.Split(input, "\"")
	if len(parts) < 2 {
		return macaroon.Macaroon{}, fmt.Errorf("invalid macaroon string")
	}

	macaroonStr := parts[1]
	return macaroon.DecodeBase64(macaroonStr)
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

			preToken, err := parsePreToken(macaroon, invoice)
			if err != nil {
				fmt.Println(err)
				return
			}

			token, err := preToken.Pay(&lightningNode)
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

func (c *TestClient) sendAuthorizationRequest(url string, token macaroon.Token) {
	fmt.Println("Sending Authorization Request...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
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

		// Store the token for later use.
		store, err := auth.NewStore("./.store")
		if err != nil {
			fmt.Println("Error creating the store:", err)
		}
		if c.tokenPath == "" {
			shareToken(store, token)
		}

		// body, _ := io.ReadAll(resp.Body)
		fmt.Println(resp.Status)
	} else {
		fmt.Println("Unexpected response status:", resp.Status)
	}
}

// getTokenPath parses the --token flag and returns its value
func getTokenPath() string {
	token := flag.String("token", "", "Path to the token file")
	flag.Parse()

	return *token
}

// shareToken stores a version of the token that can be shared with others.
func shareToken(store auth.TokenStore, token macaroon.Token) error {
	// The time at which the new macaroon will expire.
	expiryDate := time.Now().Add(5 * time.Minute)

	// Creating the restricted macaroon
	mac, err := token.Macaroon.Oven().WithThirdPartyCaveats(macaroon.NewCaveat(macaroon.ExpiryDateKey, expiryDate.Format(time.RFC3339))).Cook()
	if err != nil {
		return err
	}

	// Update the macaroon in the token.
	token.Macaroon = mac

	// Store the token.
	store.StoreToken(token.Id(), token)

	return nil
}
