package macaroon

import "fmt"

// Caveat represents a condition or restriction associated with a macaroon.
type Caveat struct {
	Key   string // The identifier or type of the caveat.
	Value string // The specific value or condition associated with the key.
}

func NewCaveat(Key string, Value string) Caveat {
	return Caveat{Key, Value}
}

func (c Caveat) String() string {
	return fmt.Sprintf("%s = %s", c.Key, c.Value)
}
