package strongoauth

import (
	"errors"
	"fmt"
	"slices"
)

func SetKnownAuthProviderIDs(v []string) {
	if knownAuthProviderIDs != nil {
		panic("knownAuthProviders already set")
	}
	knownAuthProviderIDs = v
}

var knownAuthProviderIDs []string

func ValidateAuthProviderID(v string) error {
	if v == "" {
		return errors.New("is empty string")
	}
	if slices.Contains(knownAuthProviderIDs, v) {
		return nil
	}
	return fmt.Errorf("supported auth providers=%+v, got: %s", knownAuthProviderIDs, v)
}
