package with

import (
	"github.com/strongo/slice"
)

// RolesField defines a record with a list of roles
type RolesField struct {
	Roles []string `json:"roles,omitempty" dalgo:"roles,omitempty"  firestore:"roles,omitempty"`
}

// HasRole checks if an object has a given role
func (v RolesField) HasRole(role string) bool {
	return slice.Index(v.Roles, role) >= 0
}

func (v RolesField) AddRole(role string) ( /* u dal.Update - does not make sense to return update as field unknown */ ok bool) {
	if v.HasRole(role) {
		return false
	}
	v.Roles = append(v.Roles, role)
	return true
}

// Validate returns error as soon as 1st role is not valid.
func (v RolesField) Validate() error {
	if err := ValidateSetSliceField("roles", v.Roles, true); err != nil {
		return err
	}
	return nil
}
