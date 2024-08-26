package appuser

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/validation"
	"time"
)

// WithLastLogin is a struct that contains the last login time of a user.
type WithLastLogin struct {
	LastLoginAt time.Time `json:"lastLoginAt" firestore:"lastLoginAt"`
}

// SetLastLoginAt sets the time of the last successful login.
func (l *WithLastLogin) SetLastLoginAt(time time.Time) dal.Update {
	l.LastLoginAt = time
	return dal.Update{Field: "lastLoginAt", Value: time}
}

func (l *WithLastLogin) Validate() error {
	if l.LastLoginAt.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("lastLoginAt")
	}
	return nil
}
