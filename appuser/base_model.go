package appuser

// BaseUserFields hold common properties for BaseUserData interface
type BaseUserFields struct { // former AppUserBase
	NameFields
	AccountsOfUser
	WithPreferredLocale
	WithCreatedTimestamp
}
