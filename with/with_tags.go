package with

import (
	"strings"
)

// TagsField defines a record with a list of tags
type TagsField struct {
	Tags []string `json:"tags,omitempty" dalgo:"tags,omitempty" firestore:"tags,omitempty"`
}

// Validate returns error as soon as 1st tag is not valid.
func (v TagsField) Validate() error {
	if err := ValidateSetSliceField("tags", v.Tags, false); err != nil {
		return err
	}
	return nil
}

// String returns string representation of the TagsField
func (v TagsField) String() string {
	return "tags=" + strings.Join(v.Tags, ",")
}
