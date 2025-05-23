package person

import (
	"errors"
	"fmt"
	"github.com/strongo/validation"
	"strings"
)

type NamesHolder interface {
	SetNames(names ...Name) error
	GetName(field NameField) string
}

var _ NamesHolder = (*NameFields)(nil)

// NameFields is a struct for storing names of a user or a contact
type NameFields struct {
	UserName   string `json:"userName,omitempty" dalgo:"userName,omitempty" firestore:"userName,omitempty"`
	FirstName  string `json:"firstName,omitempty" dalgo:"firstName,omitempty" firestore:"firstName,omitempty"`
	MiddleName string `json:"middleName,omitempty" dalgo:"middleName,omitempty" firestore:"middleName,omitempty"`
	LastName   string `json:"lastName,omitempty" dalgo:"lastName,omitempty" firestore:"lastName,omitempty"`
	NickName   string `json:"nickName,omitempty" dalgo:"nickName,omitempty" firestore:"nickName,omitempty"`
	ScreenName string `json:"screenName,omitempty" dalgo:"screenName,omitempty" firestore:"screenName,omitempty"`
	FullName   string `json:"fullName,omitempty" dalgo:"fullName,omitempty" firestore:"fullName,omitempty"`
}

// SetNames sets names of a user or a contact
func (v *NameFields) SetNames(names ...Name) error {
	for _, name := range names {
		switch name.Field {
		case Username:
			v.UserName = name.Value
		case FullName:
			v.FullName = name.Value
		case FirstName:
			v.FirstName = name.Value
		case MiddleName:
			v.MiddleName = name.Value
		case LastName:
			v.LastName = name.Value
		case NickName:
			v.NickName = name.Value
		case ScreenName:
			v.ScreenName = name.Value
		default:
			return fmt.Errorf("unsupported NameField field: %d", name.Field)
		}
	}
	return nil
}

// String returns string representation of a user or a contact names
func (v *NameFields) String() string {
	return fmt.Sprintf(
		`{UserName="%s", FirstName="%s", MiddleName="%s", LastName="%s", NickName="%s", FullName="%s"}`,
		v.UserName, v.FirstName, v.MiddleName, v.LastName, v.NickName, v.FullName)
}

// GetName returns a long NameField of a user or a contact
func (v *NameFields) GetName(field NameField) string {
	switch field {
	case Username:
		return v.UserName
	case FullName:
		return v.FullName
	case FirstName:
		return v.FirstName
	case MiddleName:
		return v.MiddleName
	case LastName:
		return v.LastName
	case NickName:
		return v.NickName
	case ScreenName:
		return v.ScreenName
	default:
		return ""
	}
}

// GetFullName returns a full NameField of a user or a contact
func (v *NameFields) GetFullName() string {
	if v.FullName != "" {
		return v.FullName
	}
	if v.FirstName != "" && v.NickName != "" && v.LastName != "" {
		return fmt.Sprintf("%s (%s) %s", v.FirstName, v.NickName, v.LastName)
	}
	if v.FirstName != "" && v.LastName != "" {
		return v.FirstName + " " + v.LastName
	}
	if v.FirstName != "" && v.NickName != "" {
		return fmt.Sprintf("%s (%s)", v.FirstName, v.NickName)
	}
	if v.LastName != "" && v.NickName != "" {
		return fmt.Sprintf("%s (%s)", v.LastName, v.NickName)
	}
	if v.FirstName != "" {
		return v.FirstName
	}
	if v.LastName != "" {
		return v.LastName
	}
	if v.NickName != "" {
		return v.NickName
	}
	if v.UserName != "" {
		return v.UserName
	}
	if v.ScreenName != "" {
		return v.ScreenName
	}
	return ""
}

func DeductNamesFromFullName(fullName string) (firstName, lastName string) {
	fullName = strings.TrimSpace(fullName)
	for {
		const doubleSpace = "  "
		if !strings.Contains(fullName, doubleSpace) {
			break
		}
		const singleSpace = " "
		fullName = strings.ReplaceAll(fullName, doubleSpace, singleSpace)
	}

	if names := strings.Split(fullName, " "); len(names) == 2 {
		return names[0], names[1]
	}
	return "", ""
}

func (v *NameFields) Equal(v2 *NameFields) bool {
	return v == nil && v2 == nil || v != nil && v2 != nil &&
		v.FirstName == v2.FirstName &&
		v.MiddleName == v2.MiddleName &&
		v.LastName == v2.LastName &&
		v.UserName == v2.UserName &&
		v.NickName == v2.NickName &&
		v.FullName == v2.FullName
}

// Validate validates NameField fields
func (v *NameFields) Validate() error {
	if v == nil {
		return nil
	}
	const leadingOrClosingSpaces = "leading or closing leadingOrClosingSpaces"
	if strings.TrimSpace(v.FirstName) != v.FirstName {
		return validation.NewErrBadRecordFieldValue("firstName", leadingOrClosingSpaces)
	}
	if strings.TrimSpace(v.LastName) != v.LastName {
		return validation.NewErrBadRecordFieldValue("lastName", leadingOrClosingSpaces)
	}
	if strings.TrimSpace(v.MiddleName) != v.MiddleName {
		return validation.NewErrBadRecordFieldValue("middleName", leadingOrClosingSpaces)
	}
	if strings.TrimSpace(v.FullName) != v.FullName {
		return validation.NewErrBadRecordFieldValue("fullName", leadingOrClosingSpaces)
	}
	if err := ValidateAtLeast1Name(v); err != nil {
		return err
	}
	return nil
}

// IsEmpty checks if name is empty
func (v *NameFields) IsEmpty() bool {
	return v == nil || *v == NameFields{}
}

// ValidateAtLeast1Name validates required names
func ValidateAtLeast1Name(v *NameFields) error {
	if (*v == NameFields{}) {
		return errors.New("at least 1 NameField should be specified")
	}
	//if strings.TrimSpace(v.FirstName) == "" && strings.TrimSpace(v.LastName) == "" && strings.TrimSpace(v.FullName) == "" && strings.TrimSpace(v.NickName) == "" {
	//	return validation.NewErrBadRecordFieldValue("first|last|full|nick", "at least one of names should be specified")
	//}
	return nil
}
