package with

import (
	"github.com/dal-go/dalgo/update"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// DeletedFields DTO
type DeletedFields struct {
	DeletedAt time.Time `json:"deletedAt,omitempty" dalgo:"tags,deletedAt" firestore:"deletedAt,omitempty"`
	DeletedBy string    `json:"deletedBy,omitempty"  dalgo:"tags,deletedBy" firestore:"deletedBy,omitempty"`
}

// UpdatesWhenDeleted populates update instructions for DAL when a record has been deleted
func (v *DeletedFields) UpdatesWhenDeleted() []update.Update {
	return []update.Update{
		update.ByFieldName("deletedAt", v.DeletedAt),
		update.ByFieldName("deletedBy", v.DeletedBy),
	}
}

// Validate returns error if not valid
func (v *DeletedFields) Validate() error {
	if !v.DeletedAt.IsZero() && strings.TrimSpace(v.DeletedBy) == "" {
		return validation.NewErrRecordIsMissingRequiredField("deletedBy")
	}
	if strings.TrimSpace(v.DeletedBy) != "" && v.DeletedAt.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("deletedAt")
	}
	return nil
}
