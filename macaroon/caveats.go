package macaroon

import "fmt"

type Caveat struct {
	Key   string
	Value string
}

func NewCaveat(Key string, Value string) Caveat {
	return Caveat{Key, Value}
}

func (caveat Caveat) String() string {
	return fmt.Sprintf("%s=%s", caveat.Key, caveat.Value)
}
