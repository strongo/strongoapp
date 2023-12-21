package with

import (
	"errors"
	"strings"
)

// ValidateRecordID validates record ContactID
func ValidateRecordID(id string) error { // TODO: move into github.com/strongo/validation/validate
	if v := strings.TrimSpace(id); v == "" {
		return errors.New("id is required")
	}
	if strings.TrimSpace(id) != id {
		return errors.New("id should not start or end with whitespace characters")
	}
	if strings.ContainsAny(id, " \t\r") {
		return errors.New("must not contain whitespace characters")
	}
	return nil
}
