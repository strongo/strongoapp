package appuser

import "time"

// BelongsToUser should be implemented by any struct that belongs to a single user
type BelongsToUser interface {
	GetAppUserID() (appUserID string)
	SetAppUserID(appUserID string)
	GetCreatedTime() time.Time
	GetUpdatedTime() time.Time
	SetUpdatedTime(time.Time) error
}
