package strongo

type AppUser interface {
	SetPreferredLocale(code5 string) error
	PreferredLocale() string

	SetNames(firs, last, user string)
}
