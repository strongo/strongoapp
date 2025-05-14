package with

import (
	"fmt"
	"github.com/strongo/validation"
	"strings"
)

const PhonesField = "phones"

type Phones struct {
	Phones map[string]PhoneProps `json:"phones,omitempty" firestore:"phones,omitempty"`
}

func (v Phones) Validate() error {
	for k, p := range v.Phones {
		if strings.TrimSpace(k) == "" {
			return validation.NewErrBadRecordFieldValue(PhonesField, "phone key is empty")
		}
		if trimmedKey := strings.TrimSpace(k); trimmedKey == "" {
			return validation.NewErrBadRecordFieldValue(PhonesField+fmt.Sprintf("[%s]", k), "phone key is empty")
		} else if k != trimmedKey {
			return validation.NewErrBadRecordFieldValue(PhonesField+fmt.Sprintf("[%s]", k), "phone key has leading or trailing spaces")
		}
		if err := p.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(PhonesField+fmt.Sprintf("[%s]", k), err.Error())
		}
	}
	return nil
}

type PhoneProps struct {
	CreatedFields
	TagsField
	Type     string `json:"type" firestore:"type"`
	Title    string `json:"title,omitempty" firestore:"title,omitempty"`
	Verified bool   `json:"verified,omitempty" firestore:"verified,omitempty"`
	Note     string `json:"note,omitempty" firestore:"note,omitempty"`
}

func (v PhoneProps) Validate() error {
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	if err := v.TagsField.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.Type) == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
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
	return nil
}
