package person

type NameField int

const (
	Username NameField = iota
	FirstName
	MiddleName
	LastName
	NickName
	FullName
)

// Name is a struct for setting names of a user or a contact
type Name struct {
	Field NameField
	Value string
}
