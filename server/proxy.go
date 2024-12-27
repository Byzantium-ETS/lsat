package main

import (
	"fmt"
	"lsat/auth"
	"lsat/macaroon"
	"lsat/secrets"
	"lsat/service"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lightningnetwork/lnd/lntypes"
)

const (
	macaroonHeader    = "L402"
	defaultService    = "image"
	authFailedMessage = "Authentication failed!"
)

// LSATProxyServer is a struct that contains the necessary information to handle service requests.
type LSATProxyServer struct {
	*auth.Minter
}

// Handle the minting of a new token.
func (h *LSATProxyServer) HandleMint(c *gin.Context) {
	serviceName := c.Param("service")
	// Parse the service name.
	serviceId, err := service.ParseServiceID(serviceName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mint a new token.
	uid := secrets.NewUserId()
	pretoken, err := h.Minter.MintToken(uid, serviceId)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Set the WWW-Authenticate header.
	mac := pretoken.Macaroon
	authHeader := fmt.Sprintf("%s macaroon=\"%s\", invoice=\"%s\"", macaroonHeader, mac, pretoken.InvoiceResponse.Invoice)
	c.Header("WWW-Authenticate", authHeader)
	c.JSON(http.StatusPaymentRequired, gin.H{"error": "Payment Required"})
}

// Handle the authorization of a token.
func (h *LSATProxyServer) HandleToken(c *gin.Context) {
	// Get service ID from the request
	serviceName := c.Param("service")
	serviceID, err := service.ParseServiceID(serviceName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the Authorization header from the request
	authHeader := c.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != macaroonHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
		return
	}

	credentials := strings.Split(parts[1], ":")
	if len(credentials) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	mac, err := macaroon.DecodeBase64(credentials[0])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	preimage, err := lntypes.MakePreimageFromStr(credentials[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token := macaroon.Token{
		Macaroon: mac,
		Preimage: preimage,
	}

	// Check if the token is valid.
	err = h.Minter.AuthToken(&token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Execute callbacks for this service
	if service, err := h.Minter.ServiceManager().GetService(serviceID); err == nil {
		if service.Callback != nil {
			if err := service.Callback(c); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Run the service.
func (h *LSATProxyServer) Run() {
	// Initialize the Gin router.
	router := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Allow all origins
	config.AllowMethods = []string{"GET", "PUT", "POST", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"WWW-Authenticate",
	}
	config.ExposeHeaders = []string{
		"WWW-Authenticate", // Important to expose this header for LSAT
	}

	// Use CORS middleware
	router.Use(cors.New(config))

	// Define the routes.
	router.PUT("/service/:service", h.HandleMint)
	router.GET("/service/:service", h.HandleToken)

	// Start the server.
	port := getEnv("PORT", "8080")
	router.Run("localhost:" + port)
}

// Get the value of an environment variable or a default value.
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
