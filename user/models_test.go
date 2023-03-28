package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
		if changed := accounts.AddAccount(Account{Provider: "email", ID: "test@example.com"}); !changed {
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
		if changed := accounts.AddAccount(Account{Provider: "telegram", ID: "123456"}); !changed {
			t.Error("Shoud return changed=True")
		}
	})
}

func TestAccountsOfUser_RemoveAccount(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		accounts := AccountsOfUser{}
		assert.False(t, accounts.RemoveAccount(Account{Provider: "email", ID: "test@example.com"}))
	})
	t.Run("non_empty", func(t *testing.T) {
		accounts := AccountsOfUser{
			Accounts: []string{
				"email::u1@example.com",
				"email::u2@example.com",
				"email::u3@example.com",
			},
		}
		assert.False(t, accounts.RemoveAccount(Account{Provider: "email", ID: "test@example.com"}))
		assert.True(t, accounts.RemoveAccount(Account{Provider: "email", ID: "u2@example.com"}))
	})
}
