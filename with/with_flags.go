package with

import (
	"fmt"
	"github.com/strongo/validation"
	"strings"
)

// FlagsField defines a record with a list of flags
type FlagsField struct {
	Flags []string `json:"flags,omitempty" dalgo:"flags,omitempty" firestore:"flags,omitempty"`
}

// Validate returns error as soon as 1st flag is not valid.
func (v FlagsField) Validate() error {
	for i, flag := range v.Flags {
		if strings.TrimSpace(flag) == "" {
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("flags[%v]", i))
		}
	}
	return nil
}

// String returns string representation of the TagsField
func (v FlagsField) String() string {
	return "flags=" + strings.Join(v.Flags, ",")
}
