package strongomodels

import (
	"errors"
	"time"
)

type UpdateTimeGetter interface {
	GetUpdatedTime() time.Time
}

var _ UpdateTimeGetter = (*WithUpdatedTimestamp)(nil)

// WithUpdatedTimestamp is a struct that implements UpdateTimeGetter
type WithUpdatedTimestamp struct {
	DtUpdated time.Time `json:"dtUpdated,omitempty" dalgo:"dtUpdated,omitempty" firestore:"dtUpdated,omitempty" `
}

// GetUpdatedTime returns the time the record was last updated
func (u *WithUpdatedTimestamp) GetUpdatedTime() time.Time {
	return u.DtUpdated
}

func (u *WithUpdatedTimestamp) SetUpdatedTime(t time.Time) error {
	if t.IsZero() {
		return errors.New("passed update time is zero")
	}
	u.DtUpdated = t
	return nil
}
