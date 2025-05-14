package with

import (
	"errors"
	"fmt"
	"github.com/strongo/validation"
	"net/mail"
	"strings"
)

const EmailsField = "emails"

type Emails struct {
	Emails map[string]EmailProps `json:"emails,omitempty" firestore:"emails,omitempty"`
}

func (v *Emails) Validate() error {
	var errs []error

	hasPrimaryEmail := false
	for k, p := range v.Emails {
		if trimmed := strings.TrimSpace(k); trimmed == "" {
			errs = append(errs, validation.NewErrBadRequestFieldValue(EmailsField, "email key is empty"))
			continue
		} else if k != trimmed {
			errs = append(errs, validation.NewErrBadRecordFieldValue(EmailsField+fmt.Sprintf("[%s]", k), "email key has leading or trailing spaces"))
			continue
		}
		if _, err := mail.ParseAddress(k); err != nil {
			errs = append(errs, validation.NewErrBadRecordFieldValue(EmailsField+fmt.Sprintf("[%s]", k), "invalid email address: "+err.Error()))
			continue
		}
		if err := p.Validate(); err != nil {
			errs = append(errs, validation.NewErrBadRecordFieldValue(EmailsField+fmt.Sprintf("[%s]", k), "invalid properties: "+err.Error()))
			continue
		}
		if p.OriginalEmail == k {
			errs = append(errs, validation.NewErrBadRecordFieldValue("originalEmail", "same as email key"))
		}
		if p.Type == EmailTypePrimary {
			if hasPrimaryEmail {
				errs = append(errs, validation.NewErrBadRecordFieldValue(EmailsField, "multiple primary emails"))
			} else {
				hasPrimaryEmail = true
			}
		}
	}
	if len(errs) == 1 {
		return errs[0]
	} else if len(errs) > 1 {
		return errors.Join(errs...)
	}
	return nil
}

const EmailTypePrimary = "primary"

type EmailProps struct {
	CreatedFields
	TagsField
	Type          string `json:"type" firestore:"type"`
	AuthProvider  string `json:"authProvider" dalgo:"authProvider" firestore:"authProvider"` // E.g. Email, Google, Facebook, etc.
	Title         string `json:"title,omitempty" firestore:"title,omitempty"`
	OriginalEmail string `json:"originalEmail,omitempty" firestore:"originalEmail,omitempty"`
	Verified      bool   `json:"verified,omitempty" firestore:"verified,omitempty"`
	Note          string `json:"note,omitempty" firestore:"note,omitempty"`
}

func (v *EmailProps) Validate() error {
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	if err := v.TagsField.Validate(); err != nil {
		return err
	}
	switch v.Type {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	case EmailTypePrimary, "personal", "work":
	// Are known values
	default:
		return validation.NewErrBadRecordFieldValue("type", "unknown value: "+v.Type)
	}
	if v.OriginalEmail != "" {
		if trimmed := strings.TrimSpace(v.OriginalEmail); trimmed == "" {
			if v.OriginalEmail != "" {
				return fmt.Errorf("original email has spaces but is empty")
			}
		} else if v.OriginalEmail != trimmed {
			return fmt.Errorf("original email has leading or trailing spaces")
		}
	}
	if v.Title != "" {
		if trimmed := strings.TrimSpace(v.Title); trimmed == "" {
			if v.Title != "" {
				return fmt.Errorf("title has spaces but is empty")
			}
		} else if v.Title != trimmed {
			return fmt.Errorf("title has leading or trailing spaces")
		}
	}
	if v.Note != "" {
		if trimmed := strings.TrimSpace(v.Note); trimmed == "" {
			if v.Note != "" {
				return fmt.Errorf("note has spaces but is empty")
			}
		} else if v.Note != trimmed {
			return fmt.Errorf("note has leading or trailing spaces")
		}
	}
	return nil
}
