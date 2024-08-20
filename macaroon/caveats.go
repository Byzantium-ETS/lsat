package macaroon

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	ServiceKey     string = "service"
	ExpiryDateKey  string = "expiry_date"
	PaymentHashKey string = "payment_hash"
)

// Caveat represents a condition or restriction associated with a macaroon.
type Caveat struct {
	Key   string // The identifier or type of the caveat.
	Value string // The specific value or condition associated with the key.
}

// A new caveat.
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
	parts := strings.Split(string(data), " = ")

	key := parts[0][1:len(parts[0])]
	value := parts[1][:len(parts[1])-1]

	*caveat = NewCaveat(key, value)

	return nil
}

// ValueIterator is a helper struct to iterate over the values of a specific key in a sequence of caveats.
type ValueIterator struct {
	key     string
	caveats []Caveat
}

// Create a new servicce iterator.
func NewIterator(key string, caveats []Caveat) ValueIterator {
	return ValueIterator{key, caveats}
}

// HasNext checks if there are more caveats with the specified key.
func (vi *ValueIterator) HasNext() bool {
	for i, caveat := range vi.caveats {
		if caveat.Key == vi.key {
			vi.caveats = vi.caveats[i:]
			return true
		}
	}
	return false
}

// Next returns the value of the next caveat with the specified key.
func (vi *ValueIterator) Next() string {
	for i, caveat := range vi.caveats {
		if caveat.Key == vi.key {
			// Extract the value and remove the current caveat from the slice.
			value := caveat.Value
			vi.caveats = vi.caveats[i+1:]
			return value
		}
	}
	return ""
}
