package appuser

import (
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/strongomodels"
)

var _ BaseUserData = (*BaseUserFields)(nil)

// BaseUserFields provides a base implementation of BaseUserData interface.
type BaseUserFields struct { // former AppUserBase
	person.NameFields
	strongomodels.WithCreatedTimestamp
	strongomodels.WithUpdatedTimestamp
	strongomodels.WithPreferredLocale
	AccountsOfUser // TODO: Reconsider if this should be part of base implementation, if yes extend BaseUserData interface
}

func (v BaseUserFields) Validate() error {
	if er := v.NameFields.Validate(); er != nil {
		return er
	}
	return nil
}
