package service

import (
	"fmt"
	"lsat/macaroon"
	"strings"
	"time"
)

// Condition is a condition that must be satisfied by a set of caveats.
type Condition interface {
	// Satisfy checks if the set of caveats satisfies the condition.
	Satisfy([]macaroon.Caveat) error
}

// Timeout is a condition that checks if the expiry date of a service is valid.
type Timeout struct{}

func (e Timeout) Satisfy(caveats []macaroon.Caveat) error {
	now := time.Now()
	var previousExpiry time.Time
	iter := macaroon.NewIterator(macaroon.ExpiryDateKey, caveats)

	for iter.HasNext() {
		expiryTime := iter.Next()
		// Parse the value of the time caveat as a time.Time.
		expiry, err := time.Parse(time.RFC3339, expiryTime)

		// If there is an error parsing the time, return the error.
		if err != nil {
			return err
		}

		// Each following expiry_date should be more strict or before the previous expiry date.
		if now.After(expiry) && expiry.After(previousExpiry) {
			return fmt.Errorf("the %s is passed at %s", macaroon.ExpiryDateKey, expiry)
		}

		// Update previousExpiry to the current expiry.
		previousExpiry = expiry
	}

	return nil
}

// Capabilities is a condition that checks if the service has the required capabilities.
type Capabilities struct{ Key string }

func (c Capabilities) Satisfy(caveats []macaroon.Caveat) error {
	var previousCapabilities []string

	iter := macaroon.NewIterator(c.Key, caveats)

	for iter.HasNext() {
		value := iter.Next()
		currentCapabilities := strings.Split(value, ", ")

		fmt.Println(currentCapabilities)

		// If there are previous capabilities, check that the current capabilities are a subset of them.
		if len(previousCapabilities) > 0 && !isSubset(currentCapabilities, previousCapabilities) {
			return fmt.Errorf("capabilities %v are not a subset of the previous ones %v", currentCapabilities, previousCapabilities)
		}

		// Update previousCapabilities to the current capabilities.
		previousCapabilities = currentCapabilities
	}

	return nil
}

// isSubset checks if all elements of subset are in the set.
func isSubset(subset, set []string) bool {
	setMap := make(map[string]struct{}, len(set))

	// Populate the setMap with the elements of the set.
	for _, v := range set {
		setMap[v] = struct{}{}
	}

	// Check if each element in subset is in the setMap.
	for _, v := range subset {
		if _, exists := setMap[v]; !exists {
			return false
		}
	}
	return true
}

// UniqueKey is a condition that checks if the key is unique.
type UniqueKey struct{ Key string }

func (k UniqueKey) Satisfy(caveats []macaroon.Caveat) error {
	iter := macaroon.NewIterator(k.Key, caveats)
	iter.Next()
	if iter.HasNext() {
		return fmt.Errorf("the %s should be unique", k.Key)
	}
	return nil
}
