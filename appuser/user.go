package appuser

import (
	"time"
)

// BaseUserData defines app user interface
type BaseUserData interface {
	NamesSetter
	SetPreferredLocale(code5 string) error
	GetPreferredLocale() string
	GetCreatedTime() time.Time
}

var _ BaseUserData = (*BaseUserFields)(nil)
