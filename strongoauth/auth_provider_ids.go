package strongoauth

import (
	"errors"
	"fmt"
	"slices"
)

func SetKnownAuthProviderIDs(v []string) {
	knownAuthProviderIDs = v
}

var knownAuthProviderIDs = []string{
	"password",      // Email/Password
	"emailLink",     // Email Link (Passwordless)
	"phone",         // Phone
	"google.com",    // Google
	"facebook.com",  // Facebook
	"twitter.com",   // Twitter
	"github.com",    // GitHub
	"microsoft.com", // Microsoft
	"apple.com",     // Apple
	"yahoo.com",     // Yahoo
	"telegram",      // Telegram
	"custom",        // Custom Token
	//"playgames.google.com", // Play Games Services (Android)
}

func ValidateAuthProviderID(v string) error {
	if v == "" {
		return errors.New("is empty string")
	}
	if slices.Contains(knownAuthProviderIDs, v) {
		return nil
	}
	return fmt.Errorf("supported auth providers=%+v, got: %s", knownAuthProviderIDs, v)
}
