package with

import (
	"fmt"
	"time"
)

const Iso8601DateLayout = "2006-01-02"

// ValidateDateString checks if a string is in valid ISO "YYYY-MM-DD" format
func ValidateDateString(s string) (date time.Time, err error) {
	if l := len(s); l != 10 {
		return date, fmt.Errorf("should be in YYYY-MM-DD format, got %d chars", l)
	}
	if date, err = time.Parse(Iso8601DateLayout, s); err != nil {
		return date, err
	}
	return date, nil
}
