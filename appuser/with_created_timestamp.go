package appuser

import "time"

type CreatedTimeGetter interface {
	GetCreatedTime() time.Time
}

var _ CreatedTimeGetter = (*WithCreatedTimestamp)(nil)

type WithCreatedTimestamp struct {
	DtCreated time.Time `json:"dtCreated,omitempty" firestore:"dtCreated,omitempty" dalgo:"dtCreated,omitempty"`
}

func (u *WithCreatedTimestamp) GetCreatedTime() time.Time {
	return u.DtCreated
}
