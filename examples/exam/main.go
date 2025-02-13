package main

import (
	"errors"
	"lsat/auth"
	"lsat/macaroon"
	"lsat/mock"
	"lsat/proxy"
	"lsat/secrets"
	"lsat/service"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Question struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

type Test struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
}

var (
	secretStore = secrets.NewSecretFactory()
	challenger  = mock.NewChallenger()
	db          = make(map[string]Test, 5)
)

func main() {
	config := service.NewConfig(
		service.Service{
			Name:  "test",
			Tier:  service.BaseTier,
			Price: 1000,
			FirstPartyCaveats: []service.Caveat{
				service.Expire{Delay: time.Hour * 24},
				service.GenerateID{Name: "test_id"},
				macaroon.NewCaveat("permissions", "r,w"),
			},
			Conditions: []service.Condition{
				service.Expire{},
				service.Capabilities{Key: "permissions"},
			},
			Post: func(c any) error {
				ctx := c.(*gin.Context)

				var test Test
				if err := ctx.BindJSON(&test); err != nil {
					return err
				}

				// Get the macaroon from the Authorization header
				parts := strings.Split(ctx.GetHeader("Authorization"), " ")
				credentials := strings.Split(parts[1], ":")
				mac, _ := macaroon.DecodeBase64(credentials[0])

				// Check for which test the user is authorized
				iter := mac.GetValue("test_id")
				test_id := iter.Next()
				db[test_id] = test

				// Restrict the macaroon to only read permissions
				mac, _ = mac.Oven().WithThirdPartyCaveats(
					macaroon.NewCaveat("permissions", "r"),
				).Bake()

				ctx.JSON(http.StatusOK, mac.ToJSON())
				return nil
			},
			Get: func(c any) error {
				ctx := c.(*gin.Context)

				// Get the macaroon from the Authorization header
				parts := strings.Split(ctx.GetHeader("Authorization"), " ")
				credentials := strings.Split(parts[1], ":")
				mac, _ := macaroon.DecodeBase64(credentials[0])

				// Check for which test the user is authorized
				iter := mac.GetValue("test_id")
				test_id := iter.Next()
				test, ok := db[test_id]

				if !ok {
					return errors.New("Test not found")
				}

				ctx.JSON(http.StatusOK, test)
				return nil
			},
		},
	)
	minter := auth.NewMinter(config, secretStore, challenger)
	router := proxy.L402ProxyServer{Minter: &minter}

	router.Run()
}
