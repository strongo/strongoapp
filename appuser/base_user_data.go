package appuser

import (
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
)

// BaseUserData defines base app user interface to standardize how plugins & frameworks
// can work with a custom app user record.
// The easiest way to implement this interface is to embed BaseUserFields struct into your app user record struct.
type BaseUserData interface {
	person.NamesHolder
	with.PreferredLocaleHolder
	with.CreatedTimeGetter
	with.UpdatedTimeGetter
	with.UpdateTimeSetter
}
