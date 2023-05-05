package user

import "time"

type BelongsToUser interface {
	GetAppUserID() string
	SetAppUserID(appUserID string)
}

type BelongsToUserWithIntID interface {
	BelongsToUser
	GetAppUserIntID() int64
	SetAppUserIntID(appUserID int64)
}

type CreatedTimesSetter interface {
	SetCreatedTime(time.Time)
}

type UpdatedTimeSetter interface {
	SetUpdatedTime(time.Time)
}
