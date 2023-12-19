package appuser

import "fmt"

type name int

const (
	Username name = iota
	FirstName
	MiddleName
	LastName
	NickName
	FullName
)

type Name struct {
	Field name
	Value string
}

type NamesSetter interface {
	SetNames(names ...Name)
}

var _ NamesSetter = (*NameFields)(nil)

type NameFields struct {
	UserName   string `json:"userName,omitempty" dalgo:"userName,omitempty" firestore:"userName,omitempty"`
	FirstName  string `json:"firstName,omitempty" dalgo:"firstName,omitempty" firestore:"firstName,omitempty"`
	MiddleName string `json:"middleName,omitempty" dalgo:"middleName,omitempty" firestore:"middleName,omitempty"`
	LastName   string `json:"lastName,omitempty" dalgo:"lastName,omitempty" firestore:"lastName,omitempty"`
	NickName   string `json:"nickName,omitempty" dalgo:"nickName,omitempty" firestore:"nickName,omitempty"`
	FullName   string `json:"fullName,omitempty" dalgo:"fullName,omitempty" firestore:"fullName,omitempty"`
}

func (u *NameFields) SetNames(names ...Name) {
	for _, name := range names {
		switch name.Field {
		case Username:
			u.UserName = name.Value
		case FirstName:
			u.FirstName = name.Value
		case MiddleName:
			u.MiddleName = name.Value
		case LastName:
			u.LastName = name.Value
		case NickName:
			u.NickName = name.Value
		case FullName:
			u.FullName = name.Value
		default:
			panic(fmt.Sprintf("unsupported name field: %d", name.Field))
		}
	}
}

func (u *NameFields) String() string {
	return fmt.Sprintf(`{UserName="%s", FirstName="%s", LastName="%s"}`, u.UserName, u.FirstName, u.LastName)
}

func (u *NameFields) GetName(field name) string {
	switch field {
	case FirstName:
		return u.FirstName
	case MiddleName:
		return u.MiddleName
	case LastName:
		return u.LastName
	case Username:
		return u.UserName
	case NickName:
		return u.NickName
	case FullName:
		return u.GetFullName()
	default:
		return ""
	}
}

func (u *NameFields) GetFullName() string {
	if u.FullName != "" {
		return u.FullName
	}
	if u.FirstName != "" && u.NickName != "" && u.LastName != "" {
		return fmt.Sprintf("%s (%s) %s", u.FirstName, u.NickName, u.LastName)
	}
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" && u.NickName != "" {
		return fmt.Sprintf("%s (%s)", u.FirstName, u.NickName)
	}
	if u.LastName != "" && u.NickName != "" {
		return fmt.Sprintf("%s (%s)", u.LastName, u.NickName)
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.LastName != "" {
		return u.LastName
	}
	if u.NickName != "" {
		return u.NickName
	}
	if u.UserName != "" {
		return u.UserName
	}
	return ""
}
