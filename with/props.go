package with

import (
	"fmt"
	"github.com/strongo/validation"
	"strings"
)

const CommunicationChannelTypePersonal = "personal"
const CommunicationChannelTypeWork = "work"

func validateCommunicationChannelsField(name string, value map[string]*CommunicationChannelProps, extraValidation func(k string, v *CommunicationChannelProps) error) error {
	hasPrimary := false
	for k, p := range value {
		if strings.TrimSpace(k) == "" {
			return validation.NewErrBadRecordFieldValue(name, "phone key is empty")
		}
		if trimmedKey := strings.TrimSpace(k); trimmedKey == "" {
			return validation.NewErrBadRecordFieldValue(name+fmt.Sprintf("[%s]", k), "key is empty")
		} else if k != trimmedKey {
			return validation.NewErrBadRecordFieldValue(name+fmt.Sprintf("[%s]", k), "key has leading or trailing spaces")
		}
		if err := p.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(name+fmt.Sprintf("[%s]", k), err.Error())
		}
		if p.IsPrimary {
			if hasPrimary {
				return validation.NewErrBadRecordFieldValue(EmailsFieldName, "multiple primary emails")
			} else {
				hasPrimary = true
			}
		}
		if p.Original == k {
			return validation.NewErrBadRecordFieldValue("original", "same as key")
		}
		if extraValidation != nil {
			if err := extraValidation(k, p); err != nil {
				return err
			}
		}
	}
	return nil
}

type CommunicationChannelProps struct {
	CreatedFields
	TagsField
	Original     string `json:"original,omitempty" firestore:"original,omitempty"`
	AuthProvider string `json:"authProvider,omitempty" dalgo:"authProvider,omitempty" firestore:"authProvider,omitempty"` // E.g. Email, Google, Facebook, etc.
	IsPrimary    bool   `json:"isPrimary,omitempty" firestore:"isPrimary,omitempty"`
	IsVerified   bool   `json:"isVerified,omitempty" firestore:"IsVerified,omitempty"`
	Type         string `json:"type,omitempty" firestore:"type,omitempty"`
	Title        string `json:"title,omitempty" firestore:"title,omitempty"`
	Verified     bool   `json:"verified,omitempty" firestore:"verified,omitempty"`
	Note         string `json:"note,omitempty" firestore:"note,omitempty"`
}

func (v CommunicationChannelProps) Validate() error {
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	if err := v.TagsField.Validate(); err != nil {
		return err
	}
	if v.Original != "" {
		if trimmed := strings.TrimSpace(v.Original); trimmed == "" {
			if v.Original != "" {
				return fmt.Errorf("original email has spaces but is empty")
			}
		} else if v.Original != trimmed {
			return fmt.Errorf("original email has leading or trailing spaces")
		}
	}
	if s := strings.TrimSpace(v.Type); s != v.Type {
		return validation.NewErrBadRecordFieldValue("type", "has leading or trailing spaces")
	}
	if s := strings.TrimSpace(v.Title); s != v.Title {
		return fmt.Errorf("title has leading or trailing spaces")
	}
	if s := strings.TrimSpace(v.Note); s != v.Note {
		return fmt.Errorf("note has leading or trailing spaces")
	}
	return nil
}
