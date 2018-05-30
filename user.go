package strongo

import "github.com/strongo/app/user"

// AppUser defines app user interface
type AppUser interface {
	SetPreferredLocale(code5 string) error
	GetPreferredLocale() string

	SetNames(first, last, user string)
	//GetCurrencies() []string
}

// AppUserBase hold common properties for AppUser interface
type AppUserBase struct {
	user.AccountsOfUser
	Locale    string `datastore:",noindex"`
	FirstName string `datastore:",noindex"`
	LastName  string `datastore:",noindex"`
	UserName  string `datastore:",noindex"`
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

// SetNames sets names
func (u *AppUserBase) SetNames(first, last, user string) {
	u.FirstName = first
	u.LastName = last
	u.UserName = user
}
