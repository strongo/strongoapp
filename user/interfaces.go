package user

import "time"

type BelongsToUser interface {
	GetAppUserID() interface{}
}

type BelongsToUserWithIntID interface {
	BelongsToUser
	GetAppUserIntID() int64
	SetAppUserIntID(appUserID int64)
}

type BelongsToUserWithStrID interface {
	BelongsToUser
	GetAppUserStrID() string
	SetAppUserStrID(appUserID string)
}


type CreatedTimesSetter interface {
	SetCreatedTime(time.Time)
}

type UpdatedTimeSetter interface {
	SetUpdatedTime(time.Time)
}