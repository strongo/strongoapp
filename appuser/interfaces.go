package appuser

import "time"

// BelongsToUser should be implemented by any struct that belongs to a single user
type BelongsToUser interface {
	GetAppUserID() (appUserID string)
	SetAppUserID(appUserID string)
	GetCreatedTime() time.Time
	GetUpdatedTime() time.Time
	SetUpdatedTime(time.Time)
}

//// BelongsToUserWithIntID is deprecated. Remove once OwnedByUserWithID.AppUserIntID is removed.
//type BelongsToUserWithIntID interface {
//	BelongsToUser
//	GetAppUserIntID() (appUserID int64)
//	SetAppUserIntID(appUserID int64)
//}
