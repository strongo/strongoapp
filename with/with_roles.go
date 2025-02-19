package with

import (
	"github.com/dal-go/dalgo/update"
	"github.com/strongo/slice"
	"slices"
)

// RolesField defines a record with a list of roles
type RolesField struct {
	Roles []string `json:"roles,omitempty" dalgo:"roles,omitempty"  firestore:"roles,omitempty"`
}

// HasRole checks if an object has a given role
func (v *RolesField) HasRole(role string) bool {
	return slices.Index(v.Roles, role) >= 0
}

func (v *RolesField) AddRole(role string) ( /* u update.Update - does not make sense to return update as field unknown */ ok bool) {
	if v.HasRole(role) {
		return false
	}
	v.Roles = append(v.Roles, role)
	return true
}

// RemoveRole removes a role from the list of roles, return true if the role was removed, false if the role was not found
func (v *RolesField) RemoveRole(role string) (updates []update.Update) {
	var removedCount int
	v.Roles, removedCount = slice.RemoveInPlace(v.Roles, func(item string) bool {
		return item == role
	})
	if removedCount > 0 {
		updates = []update.Update{update.ByFieldName("roles", v.Roles)}
	}
	return
}

// Validate returns error as soon as 1st role is not valid.
func (v *RolesField) Validate() error {
	if err := ValidateSetSliceField("roles", v.Roles, true); err != nil {
		return err
	}
	return nil
}
