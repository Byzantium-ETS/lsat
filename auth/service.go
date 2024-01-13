package auth

import (
	"context"
	"lsat/macaroon"
)

// ServiceLimiter is an interface defining methods for managing services and their capabilities.
type ServiceLimiter interface {
	// Services retrieves information about services with the provided names.
	Services(context.Context, ...string) ([]macaroon.Service, error)

	// Capabilities retrieves the capabilities associated with the provided services.
	// It returns a list of caveats representing the capabilities of the services.
	Capabilities(context.Context, ...macaroon.Service) ([]macaroon.Caveat, error)

	// Sign apply a seal on the macaroon that is used to determined if it was authenticated.
	//
	// The strategy is left to the service.
	//
	// TO-DO: Define a spec the service can use to sign the token.
	Sign(macaroon.Macaroon) (macaroon.Macaroon, error)
}
