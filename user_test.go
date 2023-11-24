package strongo

import (
	"testing"
)

func TestUserNames_SetNames(t *testing.T) {
	type args struct {
		names Names
	}
	strPtr := func(s string) *string {
		return &s
	}
	tests := []struct {
		name     string
		user     UserNames
		args     args
		expected UserNames
	}{
		{
			name: "first_name",
			user: UserNames{FirstName: "Jack", LastName: "Doe"},
			args: args{
				names: Names{
					FirstName: strPtr("John"),
				},
			},
			expected: UserNames{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			name: "last_name",
			user: UserNames{FirstName: "Jack", LastName: "Doe"},
			args: args{
				names: Names{
					LastName: strPtr("Jones"),
				},
			},
			expected: UserNames{
				FirstName: "Jack",
				LastName:  "Jones",
			},
		},
		{
			name: "user_name",
			user: UserNames{FirstName: "Jack", LastName: "Doe"},
			args: args{
				names: Names{
					UserName: strPtr("joker"),
				},
			},
			expected: UserNames{
				FirstName: "Jack",
				LastName:  "Doe",
				UserName:  "joker",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.user.SetNames(tt.args.names)
			if tt.user != tt.expected {
				t.Errorf("AppUserBase.SetNames() = %v, want %v", tt.user, tt.expected)
			}
		})
	}
}
