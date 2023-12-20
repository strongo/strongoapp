package appuser

import (
	"github.com/strongo/strongoapp/person"
	"testing"
)

func TestBase_SetNames(t *testing.T) {
	type args struct {
		names []person.Name
	}
	tests := []struct {
		name     string
		user     person.NameFields
		args     args
		expected person.NameFields
	}{
		{
			name: "first_name",
			user: person.NameFields{FirstName: "Jack", LastName: "Doe"},
			args: args{
				names: []person.Name{{Field: person.FirstName, Value: "John"}},
			},
			expected: person.NameFields{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			name: "last_name",
			user: person.NameFields{FirstName: "Jack", LastName: "Doe"},
			args: args{
				names: []person.Name{{Field: person.LastName, Value: "Jones"}},
			},
			expected: person.NameFields{
				FirstName: "Jack",
				LastName:  "Jones",
			},
		},
		{
			name: "user_name",
			user: person.NameFields{FirstName: "Jack", LastName: "Doe"},
			args: args{
				names: []person.Name{{Field: person.Username, Value: "joker"}},
			},
			expected: person.NameFields{
				FirstName: "Jack",
				LastName:  "Doe",
				UserName:  "joker",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.user.SetNames(tt.args.names...); err != nil {
				t.Errorf("AppUserBase.SetNames() error = %v", err)
			}
			if tt.user != tt.expected {
				t.Errorf("AppUserBase.SetNames() = %v, want %v", tt.user, tt.expected)
			}
		})
	}
}
