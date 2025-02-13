package service

import (
	"fmt"
	"lsat/macaroon"
	"time"
)

// Condition is a condition that must be satisfied by a set of caveats.
type Condition interface {
	// Satisfy checks if the set of caveats satisfies the condition.
	Satisfy(...macaroon.Caveat) error
}

// Timeout is a condition that checks if the expiry date of a service is valid.
// type Timeout struct{}

func (e Expire) Satisfy(caveats ...macaroon.Caveat) error {
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

func (c Capabilities) Satisfy(caveats ...macaroon.Caveat) error {
	var previousCapabilities string

	iter := macaroon.NewIterator(c.Key, caveats)

	for iter.HasNext() {
		currentCapabilities := iter.Next()

		// If there are previous capabilities, check that the current capabilities are a subset of them.
		if len(previousCapabilities) > 0 && !isSubstring(previousCapabilities, currentCapabilities) {
			return fmt.Errorf("capabilities %v are not a subset of the previous ones %v", currentCapabilities, previousCapabilities)
		}

		// Update previousCapabilities to the current capabilities.
		previousCapabilities = currentCapabilities
	}

	return nil
}

// UniqueKey is a condition that checks if a key is unique.
type UniqueKey struct{ Key string }

func (k UniqueKey) Satisfy(caveats ...macaroon.Caveat) error {
	iter := macaroon.NewIterator(k.Key, caveats)
	iter.Next()
	if iter.HasNext() {
		return fmt.Errorf("the %s should be unique", k.Key)
	}
	return nil
}

func (n NotBefore) Satisfy(caveats ...macaroon.Caveat) error {
	now := time.Now()
	var latestStart time.Time
	iter := macaroon.NewIterator(macaroon.NotBeforeKey, caveats)

	for iter.HasNext() {
		startTimeStr := iter.Next()
		// Parse the value of the time caveat as a time.Time
		startTime, err := time.Parse(time.RFC3339, startTimeStr)

		// If there is an error parsing the time, return the error
		if err != nil {
			return err
		}

		// Each following not_before should be less restrictive or after the previous start date
		if now.Before(startTime) || startTime.Before(latestStart) {
			return fmt.Errorf("current time %s is before the not_before date %s", now, startTime)
		}

		// Update latestStart to the current start time
		latestStart = startTime
	}

	return nil
}

// isSubset checks if the first slice is a subset of the second slice.
func isSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		matched := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}
