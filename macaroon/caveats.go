package macaroon

import "fmt"

type Caveat struct {
	Key   string // The key of the caveat.
	Value string // The value associated with the key.
}

func NewCaveat(Key string, Value string) Caveat {
	return Caveat{Key, Value}
}

func (c Caveat) String() string {
	return fmt.Sprintf("%s = %s", c.Key, c.Value)
}
