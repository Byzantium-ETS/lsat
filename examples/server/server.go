package main

import (
	"lsat/auth"
	"lsat/mock"
	"lsat/proxy"
	"lsat/secrets"
	"lsat/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// Connect to the phoenix node
	// lightningClient = phoenixd.NewPhoenixClient("http://127.0.0.1:9740", "")
	// lightningNode = phoenixd.PhoenixNode{Client: lightningClient}
	// challenger = &challenge.ChallengeFactory{
	// 	LightningNode: &lightningNode,
	// }
	secretStore = secrets.NewSecretFactory()
	challenger  = mock.NewChallenger()
)

func main() {
	config := service.NewConfig(
		service.Service{
			Name:       "image",
			Tier:       service.BaseTier,
			Duration:   time.Hour,
			Price:      100,
			Conditions: []service.Condition{service.Timeout{}},
			Callback: func(c any) error {
				ctx := c.(*gin.Context)
				ctx.Redirect(http.StatusFound, "https://picsum.photos/1000")
				return nil
			},
		},
	)
	minter := auth.NewMinter(config, secretStore, challenger)
	router := proxy.L402ProxyServer{Minter: &minter}

	router.Run()
}
