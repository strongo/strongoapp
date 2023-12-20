package strongomodels

type PreferredLocaleHolder interface {
	SetPreferredLocale(code5 string) error
	GetPreferredLocale() string
}

var _ PreferredLocaleHolder = (*WithPreferredLocale)(nil)

// WithPreferredLocale is a struct for setting preferred locale of a user or a contact
type WithPreferredLocale struct {
	PreferredLocale string `json:"preferredLocale,omitempty" dalgo:"preferredLocale,omitempty" firestore:"preferredLocale,omitempty"`
}

// SetPreferredLocale sets preferred locale
func (u *WithPreferredLocale) SetPreferredLocale(v string) error {
	u.PreferredLocale = v
	return nil
}

// GetPreferredLocale gets preferred locale
func (u *WithPreferredLocale) GetPreferredLocale() string {
	return u.PreferredLocale
}
