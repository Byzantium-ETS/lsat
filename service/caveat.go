package service

import (
	"lsat/macaroon"
	"lsat/secrets"
	"time"
)

// Caveat is an interface that represents a caveat.
type Caveat interface {
	GetKey() string   // Key returns the key of the caveat.
	GetValue() string // Value returns the value of the caveat.
}

// ToCaveat takes a Caveat and returns a macaroon.Caveat for compatibility.
func ToCaveat(c Caveat) macaroon.Caveat {
	return macaroon.Caveat{
		Key:   c.GetKey(),
		Value: c.GetValue(),
	}
}

// GenerateID is a caveat that provides a unique identifier.
type GenerateID struct {
	Name string
}

func (u GenerateID) GetKey() string {
	return u.Name
}

func (u GenerateID) GetValue() string {
	return secrets.NewSecret().String()
}

// Expire is a caveat that provides an expiry date.
type Expire struct {
	Delay time.Duration
}

func (e Expire) GetKey() string {
	return macaroon.ExpiryDateKey
}

func (e Expire) GetValue() string {
	return time.Now().Add(e.Delay).Format(time.RFC3339)
}

// NotBefore is a caveat that provides a start date before which the token is not valid
type NotBefore struct {
	Delay time.Duration
}

func (n NotBefore) GetKey() string {
	return macaroon.NotBeforeKey
}

func (n NotBefore) GetValue() string {
	return time.Now().Add(n.Delay).Format(time.RFC3339)
}
