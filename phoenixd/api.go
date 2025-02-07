// Package phoenixd provides a client for interacting with the Phoenix API.
package phoenixd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// PhoenixClient is a client for interacting with the Phoenix API.
type PhoenixClient struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

// CreateInvoiceRequest represents the request to create an invoice.
type CreateInvoiceRequest struct {
	// Description is the description of the invoice (max. 128 characters).
	Description string `json:"description,omitempty"`
	// DescriptionHash is the sha256 hash of a description.
	DescriptionHash string `json:"descriptionHash,omitempty"`
	// AmountSat is the amount requested by the invoice, in satoshi.
	AmountSat uint64 `json:"amountSat"`
	// ExternalId is an optional custom identifier to link the invoice to an external system.
	ExternalId string `json:"externalId,omitempty"`
}

// InvoiceResponse represents the response from creating an invoice.
type InvoiceResponse struct {
	// AmountSat is the amount requested by the invoice, in satoshi.
	AmountSat uint64 `json:"amountSat"`
	// PaymentHash is the payment hash of the invoice.
	PaymentHash string `json:"paymentHash"`
	// Serialized is the serialized invoice.
	Serialized string `json:"serialized"`
}

// PayInvoiceRequest represents the request to pay an invoice.
type PayInvoiceRequest struct {
	// AmountSat is an optional amount in satoshi. If unset, will pay the amount requested in the invoice.
	AmountSat uint64 `json:"amountSat,omitempty"`
	// Invoice is the BOLT11 invoice.
	Invoice string `json:"invoice"`
}

// PaymentResponse represents the response from paying an invoice.
type PaymentResponse struct {
	// RecipientAmountSat is the amount received by the recipient, in satoshi.
	RecipientAmountSat uint64 `json:"recipientAmountSat"`
	// RoutingFeeSat is the routing fee for the payment, in satoshi.
	RoutingFeeSat uint `json:"routingFeeSat"`
	// PaymentId is the internal payment ID for the payment.
	PaymentId string `json:"paymentId"`
	// PaymentHash is the payment hash of the payment.
	PaymentHash string `json:"paymentHash"`
	// PaymentPreimage is the preimage of the payment.
	PaymentPreimage string `json:"paymentPreimage"`
}

// Payment represents the details of a payment.
type Payment struct {
	// PaymentHash is the payment hash of the payment.
	PaymentHash string `json:"paymentHash"`
	// Preimage is the preimage of the payment.
	Preimage string `json:"preimage"`
	// ExternalId is the external identifier associated with the payment.
	ExternalId string `json:"externalId"`
	// Description is the description of the payment.
	Description string `json:"description"`
	// Invoice is the serialized invoice.
	Invoice string `json:"invoice"`
	// IsPaid indicates whether the payment has been paid.
	IsPaid bool `json:"isPaid"`
	// ReceivedSat is the amount received, in satoshi.
	ReceivedSat uint64 `json:"receivedSat"`
	// Fees are the fees associated with the payment, in satoshi.
	Fees uint64 `json:"fees"`
	// CompletedAt is the timestamp when the payment was completed.
	CompletedAt int64 `json:"completedAt"`
	// CreatedAt is the timestamp when the payment was created.
	CreatedAt int64 `json:"createdAt"`
}

// Channel represents the details of a channel.
type Channel struct {
	// State is the current state of the channel.
	State string `json:"state"`
	// ChannelId is the unique identifier of the channel.
	ChannelId string `json:"channelId"`
	// BalanceSat is the current balance of the channel, in satoshi.
	BalanceSat uint64 `json:"balanceSat"`
	// InboundLiquiditySat is the inbound liquidity of the channel, in satoshi.
	InboundLiquiditySat uint64 `json:"inboundLiquiditySat"`
	// CapacitySat is the total capacity of the channel, in satoshi.
	CapacitySat uint64 `json:"capacitySat"`
	// FundingTxId is the transaction ID of the funding transaction.
	FundingTxId string `json:"fundingTxId"`
}

// NodeInfo represents the information about a node.
type NodeInfo struct {
	// NodeId is the unique identifier of the node.
	NodeId string `json:"nodeId"`
	// Channels is a list of channels associated with the node.
	Channels []Channel `json:"channels"`
}

// NewPhoenixClient creates a new PhoenixClient.
func NewPhoenixClient(baseURL, apiKey string) *PhoenixClient {
	return &PhoenixClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		APIKey:     apiKey,
	}
}

// createAuthHeader creates the Basic Authentication header.
func (c *PhoenixClient) createAuthHeader() string {
	auth := ":" + c.APIKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// CreateInvoice creates a new invoice.
func (c *PhoenixClient) CreateInvoice(req *CreateInvoiceRequest) (*InvoiceResponse, error) {
	url := fmt.Sprintf("%s/createinvoice", c.BaseURL)
	formData := fmt.Sprintf("description=%s&amountSat=%d&externalId=%s", req.Description, req.AmountSat, req.ExternalId)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(formData)))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Authorization", c.createAuthHeader())

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var invoiceResponse InvoiceResponse
	if err := json.Unmarshal(body, &invoiceResponse); err != nil {
		return nil, err
	}

	return &invoiceResponse, nil
}

// PayInvoice pays a BOLT11 Lightning invoice.
func (c *PhoenixClient) PayInvoice(req *PayInvoiceRequest) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/payinvoice", c.BaseURL)
	formData := fmt.Sprintf("invoice=%s", req.Invoice)
	if req.AmountSat != 0 {
		formData += fmt.Sprintf("&amountSat=%d", req.AmountSat)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBufferString(formData))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Authorization", c.createAuthHeader())

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var paymentResponse PaymentResponse
	if err := json.Unmarshal(body, &paymentResponse); err != nil {
		return nil, err
	}

	if paymentResponse.PaymentPreimage == "" {
		return nil, errors.New(string(body))
	}

	return &paymentResponse, nil
}

// GetIncomingPayment retrieves the details of an incoming payment.
func (c *PhoenixClient) GetIncomingPayment(paymentHash string) (*Payment, error) {
	url := fmt.Sprintf("%s/payments/incoming/%s", c.BaseURL, paymentHash)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Authorization", c.createAuthHeader())

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var payment Payment
	if err := json.Unmarshal(body, &payment); err != nil {
		return nil, err
	}

	return &payment, nil
}

// GetInfo retrieves information about the node.
func (c *PhoenixClient) GetInfo() (*NodeInfo, error) {
	url := fmt.Sprintf("%s/getinfo", c.BaseURL)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Authorization", c.createAuthHeader())

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var nodeInfo NodeInfo
	if err := json.Unmarshal(body, &nodeInfo); err != nil {
		return nil, err
	}

	return &nodeInfo, nil
}
