package with

import (
	"errors"
	"github.com/dal-go/dalgo/update"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// Created is intended to be used only in WithCreatedField. For root level use WithCreatedFields instead.
type Created struct {
	At string `json:"at" dalgo:"at" firestore:"at"`
	By string `json:"by" dalgo:"at" firestore:"by"`
}

// Validate returns error if not valid
func (v *Created) Validate() error {
	var errs []error
	if strings.TrimSpace(v.At) == "" {
		errs = append(errs, validation.NewErrRecordIsMissingRequiredField("at"))
	}
	if strings.TrimSpace(v.By) == "" {
		errs = append(errs, validation.NewErrRecordIsMissingRequiredField("by"))
	}
	if len(v.At) == len(time.DateOnly) {
		// TODO: this does not feels right, temporary workaround for values passed from client in new member creation
		if _, err := time.Parse(time.DateOnly, v.At); err != nil {
			return validation.NewErrBadRecordFieldValue("at", err.Error())
		}
	} else if _, err := time.Parse(time.RFC3339, v.At); err != nil {
		return validation.NewErrBadRecordFieldValue("at", err.Error())
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// CreatedField adds a Created field to a data model
type CreatedField struct {
	Created Created `json:"created" firestore:"created"`
}

func (v *CreatedField) Validate() error {
	if err := v.Created.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("created", err.Error())
	}
	return nil
}

// CreatedFields adds CreatedAtField and CreatedByField fields to a data model
type CreatedFields struct {
	CreatedAtField
	CreatedByField
}

// UpdatesWhenCreated populates update instructions for DAL when a record has been created
func (v *CreatedFields) UpdatesWhenCreated() []update.Update {
	return append(v.UpdatesCreatedOn(), v.UpdatesCreatedBy()...)
}

// Validate returns error if not valid
func (v *CreatedFields) Validate() error {
	var errs []error
	if err := v.CreatedAtField.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err := v.CreatedByField.Validate(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
