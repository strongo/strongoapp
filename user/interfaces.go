package user

import "time"

type BelongsToUser interface {
	GetAppUserIntID() int64
	SetAppUserIntID(appUserID int64)
	SetDtCreated(time time.Time)
	SetDtUpdated(time time.Time)
}


