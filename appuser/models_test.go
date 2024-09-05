package appuser

import (
	"github.com/strongo/strongoapp/with"
	"strings"
	"testing"
	"time"
)

func TestAccountsOfUser_AddAccount(t *testing.T) {
	t.Parallel()

	verifyOutput := func(t *testing.T, accounts AccountsOfUser, count int) {
		if len(accounts.Accounts) != count {
			t.Errorf("len(accounts) != %d", count)
		}
	}

	t.Run("email", func(t *testing.T) {
		accounts := AccountsOfUser{}
		if updates := accounts.AddAccount(AccountKey{Provider: "email", ID: "test@example.com"}); len(updates) == 0 {
			t.Error("should not return any updates")
		}
		verifyOutput(t, accounts, 1)
	})

	t.Run("panic on missing app", func(t *testing.T) {
		accounts := AccountsOfUser{}
		defer func() {
			if r := recover(); r == nil {
				t.Error("should panic")
			}
			verifyOutput(t, accounts, 0)
		}()
		if updates := accounts.AddAccount(AccountKey{Provider: "facebook", ID: "123456"}); len(updates) != 1 {
			t.Errorf("should return 1 update, got %d", len(updates))
		}
	})
	t.Run("do_not_panic_on_missing_app_for_telegram", func(t *testing.T) {
		accounts := AccountsOfUser{}
		if updates := accounts.AddAccount(AccountKey{Provider: "telegram", ID: "123456"}); len(updates) != 1 {
			t.Errorf("should return 1 update, got %d", len(updates))
			verifyOutput(t, accounts, 1)
		}
	})
}

func TestAccountsOfUser_RemoveAccount(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		accounts := AccountsOfUser{}
		changed := accounts.RemoveAccount(AccountKey{Provider: "email", ID: "test@example.com"})
		if changed {
			t.Error("Should return changed=False")
		}
	})
	t.Run("non_empty", func(t *testing.T) {
		accounts := AccountsOfUser{
			Accounts: []string{
				"email::u1@example.com",
				"email::u2@example.com",
				"email::u3@example.com",
			},
		}
		if changed := accounts.RemoveAccount(AccountKey{Provider: "email", ID: "test@example.com"}); changed {
			t.Error("Should return changed=False")
		}
		if changed := accounts.RemoveAccount(AccountKey{Provider: "email", ID: "u2@example.com"}); !changed {
			t.Error("Should return changed=True")
		}
	})
}

func TestNewOwnedByUserWithID(t *testing.T) {
	type args struct {
		id      string
		created time.Time
	}
	now := time.Now()
	tests := []struct {
		name        string
		args        args
		shouldPanic bool
		want        OwnedByUserWithID
	}{
		{name: "should_pass", args: args{id: "123", created: now},
			want: OwnedByUserWithID{
				AppUserID: "123",
				CreatedFields: with.CreatedFields{
					CreatedAtField: with.CreatedAtField{CreatedAt: now},
				},
			},
		},
		{name: "empty_id", args: args{created: now}, shouldPanic: true},
		{name: "empty_created", args: args{id: "123"}, shouldPanic: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("should panic")
					}
				}()
			}
			actual := NewOwnedByUserWithID(tt.args.id, tt.args.created)
			if actual != tt.want {
				t.Errorf("NewOwnedByUserWithID(%v, %v) = %v, want %v", tt.args.id, tt.args.created, actual, tt.want)
			}
		})
	}
}

func TestNewEmailData(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want EmailData
	}{
		{
			name: "empty",
			args: args{email: ""},
			want: EmailData{},
		},
		{
			name: "not_empty",
			args: args{
				email: "   Test@eXample.com  ",
			},
			want: EmailData{
				EmailRaw:       "Test@eXample.com",
				EmailLowerCase: "test@example.com",
				EmailConfirmed: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEmailData(tt.args.email); got != tt.want {
				t.Errorf("NewEmailData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountKey_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		ua   AccountKey
		want bool
	}{
		{
			name: "empty",
			ua:   AccountKey{},
			want: true,
		},
		{
			name: "app",
			ua:   AccountKey{App: "test"},
			want: false,
		},
		{
			name: "provider",
			ua:   AccountKey{Provider: "test"},
			want: false,
		},
		{
			name: "id",
			ua:   AccountKey{ID: "test"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ua.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountKey_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ua      AccountKey
		wantErr string
	}{
		{
			name:    "no_app",
			ua:      AccountKey{Provider: "p1", ID: "id1"},
			wantErr: "",
		},
		{
			name:    "full",
			ua:      AccountKey{Provider: "p1", ID: "id1", App: "app1"},
			wantErr: "",
		},
		{
			name:    "empty",
			ua:      AccountKey{},
			wantErr: "empty struct",
		},
		{
			name:    "provider",
			ua:      AccountKey{ID: "id1"},
			wantErr: "[provider]",
		},
		{
			name:    "id",
			ua:      AccountKey{Provider: "p1"},
			wantErr: "[id]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ua.Validate()
			if err == nil && tt.wantErr != "" {
				t.Error("should return an error")
			} else if err != nil {
				if tt.wantErr == "" {
					t.Errorf("should not return an error, got: %v", err)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error message mismatch: '%s' != '%s'", err, tt.wantErr)
				}
			}
		})
	}
}
