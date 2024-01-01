package auth

import (
	"context"
	"lsat/macaroon"
)

// ServiceManager is an interface defining methods for managing services and their capabilities.
type ServiceManager interface {
	// Services retrieves information about services with the provided names.
	Services(cx context.Context, names ...string) ([]macaroon.Service, error)

	// Capabilities retrieves the capabilities associated with the provided services.
	// It returns a list of caveats representing the capabilities of the services.
	Capabilities(cx context.Context, services ...macaroon.Service) ([]macaroon.Caveat, error)

	// VerifyCaveats verifies the validity of a set of macaroon caveats.
	// It returns an error if any of the caveats are invalid or do not meet the specified criteria.
	VerifyCaveats(caveats ...macaroon.Caveat) error

	// GetResource retrieves the resource associated with the provided service.
	// The resource can be, for example, an image, file, or any data that the service wants to provide.
	GetResource(cx context.Context, macaroon macaroon.Macaroon) (Resource, error)
}

// Resource represents the data associated with a service's resource.
type Resource struct {
	// Add fields representing the details of the resource, such as content, type, etc.
	Content []byte
	Type    string
}
