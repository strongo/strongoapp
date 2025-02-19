package with

import (
	"github.com/dal-go/dalgo/update"
	"github.com/strongo/validation"
	"time"
)

type CreatedTimeGetter interface {
	GetCreatedTime() (time.Time, error)
}

var _ CreatedTimeGetter = (*CreatedAtField)(nil)

type CreatedAtField struct {
	//CreatedAt string `json:"createdAt" dalgo:"createdAt" firestore:"createdAt"`
	CreatedAt time.Time `json:"createdAt" dalgo:"createdAt" firestore:"createdAt"`
}

// GetCreatedTime returns value of CreatedAt field as time.Time parsed with RFC3339Nano layout
func (v *CreatedAtField) GetCreatedTime() (time.Time, error) {
	//return time.Parse(time.RFC3339Nano, v.CreatedAt)
	return v.CreatedAt, nil
}

// SetCreatedAt sets CreatedAtField field formatted with RFC3339Nano layout
func (v *CreatedAtField) SetCreatedAt(t time.Time) {
	//v.CreatedAt = t.Format(time.RFC3339Nano)
	v.CreatedAt = t
}

func (v *CreatedAtField) UpdatesCreatedOn() []update.Update {
	return []update.Update{
		update.ByFieldName("createdOn", v.CreatedAt),
	}
}

func (v *CreatedAtField) Validate() error {
	if v.CreatedAt.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("createdAt")
	}
	//if _, err := time.Parse(time.DateOnly, v.CreatedAt); err != nil {
	//	return validation.NewErrBadRecordFieldValue("createdAt", err.Error())
	//}
	return nil
}
