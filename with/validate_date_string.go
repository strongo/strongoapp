package with

import (
	"fmt"
	"time"
)

// ValidateDateString checks if a string is in valid ISO "YYYY-MM-DD" format
func ValidateDateString(s string) (date time.Time, err error) {
	if l := len(s); l != 10 {
		return date, fmt.Errorf("should be in YYYY-MM-DD format, got %d chars", l)
	}
	if date, err = time.Parse(time.DateOnly, s); err != nil {
		return date, err
	}
	return date, nil
}
