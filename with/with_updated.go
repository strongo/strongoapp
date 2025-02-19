package with

import (
	"errors"
	"github.com/dal-go/dalgo/update"
	"github.com/strongo/validation"
	"strings"
	"time"
)

type UpdatedTimeGetter interface {
	GetUpdatedTime() time.Time
}

type UpdateTimeSetter interface {
	SetUpdatedTime(time.Time) error
}

var _ UpdatedTimeGetter = (*UpdatedFields)(nil)

// GetUpdatedTime returns the time the record was last updated
func (u *UpdatedFields) GetUpdatedTime() time.Time {
	return u.UpdatedAt
}

// SetUpdatedTime sets UpdatedAt field to the time provided
func (u *UpdatedFields) SetUpdatedTime(t time.Time) error {
	u.UpdatedAt = t
	return nil
}

// UpdatedFields provides UpdatedAt & UpdatedBy fields
type UpdatedFields struct {
	UpdatedAt time.Time `json:"updatedAt,omitempty"  firestore:"updatedAt,omitempty"`
	UpdatedBy string    `json:"updatedBy,omitempty"  firestore:"updatedBy,omitempty"`
}

// UpdatesWhenUpdatedFieldsChanged populates update instructions for DALgo when UpdatedAt or UpdatedBy fields changed
func (v *UpdatedFields) UpdatesWhenUpdatedFieldsChanged() []update.Update {
	return []update.Update{
		update.ByFieldName("updatedAt", v.UpdatedAt),
		update.ByFieldName("updatedBy", v.UpdatedBy),
	}
}

// Validate returns error if not valid
func (v *UpdatedFields) Validate() error {
	var errs []error
	if v.UpdatedAt.IsZero() {
		errs = append(errs, validation.NewErrRecordIsMissingRequiredField("updatedAt"))
	}
	if strings.TrimSpace(v.UpdatedBy) == "" {
		errs = append(errs, validation.NewErrRecordIsMissingRequiredField("updatedBy"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
