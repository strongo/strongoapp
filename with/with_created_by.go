package with

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/validation"
	"strings"
)

type CreatedByField struct {
	CreatedBy string `json:"createdBy" dalgo:"createdBy" firestore:"createdBy"`
}

// SetCreatedBy sets CreatedByField field
func (v *CreatedByField) SetCreatedBy(createBy string) {
	v.CreatedBy = createBy
}

// GetCreatedBy returns CreatedByField
func (v *CreatedByField) GetCreatedBy() string {
	return v.CreatedBy
}

func (v *CreatedByField) UpdatesCreatedBy() []dal.Update {
	return []dal.Update{
		{Field: "createdBy", Value: v.CreatedBy},
	}
}

func (v *CreatedByField) Validate() error {
	if strings.TrimSpace(v.CreatedBy) == "" {
		return validation.NewErrRecordIsMissingRequiredField("createdBy")
	}
	return nil
}
