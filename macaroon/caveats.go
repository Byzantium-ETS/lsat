package macaroon

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Caveat represents a condition or restriction associated with a macaroon.
type Caveat struct {
	Key   string // The identifier or type of the caveat.
	Value string // The specific value or condition associated with the key.
}

func NewCaveat(Key string, Value string) Caveat {
	return Caveat{Key, Value}
}

func (caveat Caveat) String() string {
	return fmt.Sprintf("%s = %s", caveat.Key, caveat.Value)
}

func (caveat *Caveat) MarshalJSON() ([]byte, error) {
	return json.Marshal(caveat.String())
}

func (caveat *Caveat) UnmarshalJSON(data []byte) error {
	// fmt.Println(string(data))
	parts := strings.Split(string(data), " = ")

	key := parts[0][1:len(parts[0])]
	value := parts[1][:len(parts[1])-1]

	*caveat = NewCaveat(key, value)

	return nil
}
