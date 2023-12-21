package appuser

import (
	"github.com/strongo/strongoapp/with"
)

// BelongsToUser should be implemented by any struct that belongs to a single user
type BelongsToUser interface {
	GetAppUserID() (appUserID string)
	SetAppUserID(appUserID string)
	with.CreatedTimeGetter
	with.UpdatedTimeGetter
	with.UpdateTimeSetter
}
