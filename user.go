package strongo

type AppUser interface {
	SetPreferredLocale(code5 string) error
	GetPreferredLocale() string

	SetNames(first, last, user string)
	//GetCurrencies() []string
}

type AppUserBase struct {
	Locale    string `datastore:",noindex"`
	FirstName string `datastore:",noindex"`
	LastName  string `datastore:",noindex"`
	UserName  string `datastore:",noindex"`
}

var _ AppUser = (*AppUserBase)(nil)

func (u *AppUserBase) SetPreferredLocale(v string) error {
	u.Locale = v
	return nil
}

func (u *AppUserBase) GetPreferredLocale() string {
	return u.Locale
}

func (u *AppUserBase) SetNames(first, last, user string) {
	u.FirstName = first
	u.LastName = last
	u.UserName = user
}
