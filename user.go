package strongo

type AppUser interface {
	SetPreferredLocale(code5 string) error
	PreferredLocale() string

	SetNames(first, last, user string)
	GetCurrencies() []string // TODO: Remove from the interface
}
