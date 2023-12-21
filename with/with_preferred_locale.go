package with

type PreferredLocaleHolder interface {
	SetPreferredLocale(code5 string) error
	GetPreferredLocale() string
}

var _ PreferredLocaleHolder = (*PreferredLocaleField)(nil)

// PreferredLocaleField is a struct for setting preferred locale of a user or a contact
type PreferredLocaleField struct {
	PreferredLocale string `json:"preferredLocale,omitempty" dalgo:"preferredLocale,omitempty" firestore:"preferredLocale,omitempty"`
}

// SetPreferredLocale sets preferred locale
func (u *PreferredLocaleField) SetPreferredLocale(v string) error {
	u.PreferredLocale = v
	return nil
}

// GetPreferredLocale gets a preferred locale
func (u *PreferredLocaleField) GetPreferredLocale() string {
	return u.PreferredLocale
}
