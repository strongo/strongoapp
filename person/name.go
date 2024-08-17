package person

type NameField int

const (
	Username NameField = iota
	FullName
	FirstName
	MiddleName
	LastName
	NickName
	ScreenName
)

// Name is a struct for setting names of a user or a contact
type Name struct {
	Field NameField
	Value string
}
