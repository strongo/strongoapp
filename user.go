package strongo

import (
	"fmt"
	"github.com/strongo/app/user"
	"time"
)

type Names struct {
	FirstName *string
	LastName  *string
	UserName  *string
}

// AppUser defines app user interface
type AppUser interface {
	SetPreferredLocale(code5 string) error
	GetPreferredLocale() string

	SetNames(names Names)
	//GetCurrencies() []string
}

// AppUserBase hold common properties for AppUser interface
type AppUserBase struct {
	UserNames
	DtCreated time.Time
	user.AccountsOfUser
	Locale string `datastore:",noindex"`
}

func (u *AppUserBase) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.LastName != "" {
		return u.LastName
	}
	if u.UserName != "" {
		return u.UserName
	}
	return ""
}

var _ AppUser = (*AppUserBase)(nil)

// SetPreferredLocale sets preferred locale
func (u *AppUserBase) SetPreferredLocale(v string) error {
	u.Locale = v
	return nil
}

// GetPreferredLocale gets preferred locale
func (u *AppUserBase) GetPreferredLocale() string {
	return u.Locale
}

type UserNames struct {
	FirstName string `datastore:",noindex"`
	LastName  string `datastore:",noindex"`
	UserName  string `datastore:",noindex"`
}

func (u *UserNames) String() string {
	return fmt.Sprintf(`{UserName="%s", FirstName="%s", LastName="%s"}`, u.UserName, u.FirstName, u.LastName)
}

// SetNames sets names
func (u *UserNames) SetNames(names Names) {
	if names.FirstName != nil {
		u.FirstName = *names.FirstName
	}
	if names.LastName != nil {
		u.LastName = *names.LastName
	}
	if names.UserName != nil {
		u.UserName = *names.UserName
	}
}
