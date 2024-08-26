package appuser

import (
	"github.com/strongo/strongoapp/with"
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
		if changed := accounts.AddAccount(AccountKey{Provider: "email", ID: "test@example.com"}); !changed {
			t.Error("Shoud return changed=True")
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
		if changed := accounts.AddAccount(AccountKey{Provider: "facebook", ID: "123456"}); !changed {
			t.Error("Should return changed=True")
		}
	})
	t.Run("do_not_panic_on_missing_app_for_telegram", func(t *testing.T) {
		accounts := AccountsOfUser{}
		if changed := accounts.AddAccount(AccountKey{Provider: "telegram", ID: "123456"}); !changed {
			t.Error("Should return changed=True")
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
