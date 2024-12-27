package main

import (
	"lsat/auth"
	"lsat/mock"
	"lsat/secrets"
	"lsat/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Connect to the phoenix node
// var lightningClient = phoenixd.NewPhoenixClient("http://127.0.0.1:9740", "")
// var lightningNode = phoenixd.PhoenixNode{Client: lightningClient}

var (
	secretStore = secrets.NewSecretFactory()
	challenger  = mock.NewChallenger()
	// challenger = &challenge.ChallengeFactory{
	// 	LightningNode: &lightningNode,
	// }
)

func main() {
	config := service.NewConfig(
		service.Service{
			Name:       "image",
			Tier:       service.BaseTier,
			Duration:   time.Hour,
			Conditions: []service.Condition{service.Timeout{}},
			Callback: func(c any) error {
				ctx := c.(*gin.Context)
				ctx.Redirect(http.StatusFound, "https://picsum.photos/1000")
				return nil
			},
		},
	)
	minter := auth.NewMinter(config, secretStore, challenger)
	router := LSATProxyServer{Minter: &minter}

	router.Run()
}
