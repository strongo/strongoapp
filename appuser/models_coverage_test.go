package appuser

import (
	"testing"
	"time"

	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
)

func TestOwnedByUserWithID_Coverage(t *testing.T) {
	now := time.Now()

	// GetAppUserID, SetAppUserID, SetAppUserIntID
	o := OwnedByUserWithID{}
	o.SetAppUserID("user1")
	if o.GetAppUserID() != "user1" {
		t.Fatalf("expected user1, got %s", o.GetAppUserID())
	}
	o.SetAppUserIntID(999)
	if o.GetAppUserID() != "999" {
		t.Fatalf("expected 999, got %s", o.GetAppUserID())
	}

	// Validate: AppUserID is empty
	o2 := OwnedByUserWithID{}
	if err := o2.Validate(); err == nil || err.Error() != "AppUserID is required field" {
		t.Fatalf("expected AppUserID is required field error, got %v", err)
	}

	// Validate: CreatedAt is zero
	o2.AppUserID = "u1"
	if err := o2.Validate(); err == nil || err.Error() != "DtCreated.IsZero()" {
		t.Fatalf("expected DtCreated.IsZero() error, got %v", err)
	}

	// Validate: UpdatedAt is zero -> sets UpdatedAt = CreatedAt
	o2.CreatedAt = now
	o2.UpdatedAt = time.Time{}
	if err := o2.Validate(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if o2.UpdatedAt != now {
		t.Fatalf("expected UpdatedAt to equal CreatedAt")
	}

	// Validate: UpdatedAt before CreatedAt
	o2.UpdatedAt = now.Add(-time.Hour)
	if err := o2.Validate(); err == nil || err.Error() != "DtUpdated.Before(DtCreated) is true" {
		t.Fatalf("expected DtUpdated.Before(DtCreated) is true, got %v", err)
	}
}

func TestEmailData_Coverage(t *testing.T) {
	// Validate
	ed := EmailData{EmailRaw: "TEST@EXAMPLE.COM", EmailLowerCase: "test@example.com"}
	if err := ed.Validate(); err == nil {
		t.Fatal("expected error from Validate when EmailRaw is lowercase matched")
	}

	edMismatch := EmailData{EmailRaw: "foo", EmailLowerCase: "bar"}
	if err := edMismatch.Validate(); err == nil {
		t.Fatal("expected error from Validate when mismatch")
	}

	edEmpty := EmailData{}
	if err := edEmpty.Validate(); err != nil {
		t.Fatalf("expected nil error on empty EmailData, got %v", err)
	}

	// GetEmailRaw
	ed1 := EmailData{EmailRaw: "raw@test.com"}
	if ed1.GetEmailRaw() != "raw@test.com" {
		t.Fatalf("expected raw@test.com, got %s", ed1.GetEmailRaw())
	}
	ed2 := EmailData{EmailLowerCase: "lower@test.com"}
	if ed2.GetEmailRaw() != "lower@test.com" {
		t.Fatalf("expected lower@test.com, got %s", ed2.GetEmailRaw())
	}
	ed3 := EmailData{}
	if ed3.GetEmailRaw() != "" {
		t.Fatalf("expected empty string, got %s", ed3.GetEmailRaw())
	}

	// GetEmailLowerCase
	if ed2.GetEmailLowerCase() != "lower@test.com" {
		t.Fatalf("expected lower@test.com, got %s", ed2.GetEmailLowerCase())
	}

	// GetEmailConfirmed & SetEmailConfirmed
	if ed1.GetEmailConfirmed() != false {
		t.Fatal("expected false")
	}
	ed1.SetEmailConfirmed(true)
	if ed1.GetEmailConfirmed() != true {
		t.Fatal("expected true")
	}
}

func TestAccountDataBase_Coverage(t *testing.T) {
	now := time.Now()

	adb := AccountDataBase{
		OwnedByUserWithID: OwnedByUserWithID{
			AppUserID: "u1",
			CreatedFields: withCreated(now),
			UpdatedFields: withUpdated(now),
		},
		NameFields: person.NameFields{
			FirstName: "Alice",
		},
		WithLastLogin: WithLastLogin{
			LastLoginAt: now,
		},
	}

	// Success Validate
	if err := adb.Validate(); err != nil {
		t.Fatalf("expected valid AccountDataBase, got %v", err)
	}

	// GetNames
	names := adb.GetNames()
	if names.FirstName != "Alice" {
		t.Fatalf("expected Alice, got %s", names.FirstName)
	}

	// OwnedByUserWithID error
	bad1 := adb
	bad1.OwnedByUserWithID.AppUserID = ""
	if err := bad1.Validate(); err == nil {
		t.Fatal("expected error on missing AppUserID")
	}

	// NameFields error
	bad2 := adb
	bad2.NameFields = person.NameFields{FirstName: " Alice"}
	if err := bad2.Validate(); err == nil {
		t.Fatal("expected error on bad NameFields")
	}

	// WithLastLogin error
	bad3 := adb
	bad3.WithLastLogin = WithLastLogin{}
	if err := bad3.Validate(); err == nil {
		t.Fatal("expected error on zero LastLoginAt")
	}

	// EmailData error
	bad4 := adb
	bad4.EmailData = EmailData{EmailRaw: "foo", EmailLowerCase: "bar"}
	if err := bad4.Validate(); err == nil {
		t.Fatal("expected error on EmailData mismatch")
	}
}

func withCreated(t time.Time) (cf with.CreatedFields) {
	cf.SetCreatedAt(t)
	return cf
}

func withUpdated(t time.Time) (uf with.UpdatedFields) {
	_ = uf.SetUpdatedTime(t)
	return uf
}

func TestAccountKey_ValidateSpaces(t *testing.T) {
	akIDSpace := AccountKey{Provider: "google", App: "app", ID: " id "}
	if err := akIDSpace.Validate(); err == nil {
		t.Fatal("expected error for id spaces")
	}

	akProvSpace := AccountKey{Provider: " google ", App: "app", ID: "id"}
	if err := akProvSpace.Validate(); err == nil {
		t.Fatal("expected error for provider spaces")
	}

	akAppSpace := AccountKey{Provider: "google", App: " app ", ID: "id"}
	if err := akAppSpace.Validate(); err == nil {
		t.Fatal("expected error for app spaces")
	}
}

func TestParseUserAccount_Coverage(t *testing.T) {
	// Case 2 parts
	ak2, err := ParseUserAccount("email:user@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ak2.Provider != "email" || ak2.App != "" || ak2.ID != "user@example.com" {
		t.Fatalf("unexpected parsed result: %+v", ak2)
	}

	// Case default (error)
	if _, err := ParseUserAccount("singlepart"); err == nil {
		t.Fatal("expected error for 1 part")
	}
	if _, err := ParseUserAccount("a:b:c:d"); err == nil {
		t.Fatal("expected error for 4 parts")
	}
}

func TestAccountsOfUser_Coverage(t *testing.T) {
	// AddAccount Panics:
	// 1. ID empty or "0"
	assertPanic(t, "empty ID", func() {
		var a AccountsOfUser
		a.AddAccount(AccountKey{Provider: "email", ID: ""})
	})
	assertPanic(t, "zero ID", func() {
		var a AccountsOfUser
		a.AddAccount(AccountKey{Provider: "email", ID: "0"})
	})
	// 2. ID contains colon
	assertPanic(t, "colon in ID", func() {
		var a AccountsOfUser
		a.AddAccount(AccountKey{Provider: "email", ID: "id:extra"})
	})
	// 3. App is empty with invalid auth provider
	assertPanic(t, "unknown provider with empty app", func() {
		var a AccountsOfUser
		a.AddAccount(AccountKey{Provider: "unknown_provider", ID: "123"})
	})
	// 4. App contains colon
	assertPanic(t, "colon in App", func() {
		var a AccountsOfUser
		a.AddAccount(AccountKey{Provider: "telegram", App: "app:colon", ID: "123"})
	})
	// 5. Provider is empty
	assertPanic(t, "empty Provider", func() {
		var a AccountsOfUser
		a.AddAccount(AccountKey{Provider: "", App: "app1", ID: "123"})
	})
	// 6. Provider contains colon
	assertPanic(t, "colon in Provider", func() {
		var a AccountsOfUser
		a.AddAccount(AccountKey{Provider: "prov:colon", App: "app1", ID: "123"})
	})

	// Duplicate account returns without updates
	var a AccountsOfUser
	upd1 := a.AddAccount(AccountKey{Provider: "emailLink", ID: "u1@example.com"})
	if len(upd1) != 1 {
		t.Fatalf("expected 1 update, got %d", len(upd1))
	}
	upd2 := a.AddAccount(AccountKey{Provider: "emailLink", ID: "u1@example.com"})
	if len(upd2) != 0 {
		t.Fatalf("expected 0 updates on duplicate, got %d", len(upd2))
	}

	// SetBotUserID
	a.SetBotUserID("telegram", "mybot", "98765")
	if !a.HasAccount("telegram", "mybot") {
		t.Fatal("expected telegram:mybot to exist")
	}

	// HasAccount with empty app and non-empty app
	if !a.HasAccount("emailLink", "") {
		t.Fatal("expected email account to exist")
	}
	if a.HasAccount("nonexistent", "") {
		t.Fatal("expected nonexistent account to not exist")
	}

	// Deprecated Has* panics
	assertPanic(t, "HasTelegramAccount deprecated", func() {
		a.HasTelegramAccount()
	})
	assertPanic(t, "HasGoogleAccount deprecated", func() {
		a.HasGoogleAccount()
	})

	// GetTelegramUserIDs
	var tgAccounts AccountsOfUser
	tgAccounts.Accounts = []string{"telegram::111", "telegram::222", "email::u@test.com"}
	tgIDs := tgAccounts.GetTelegramUserIDs()
	if len(tgIDs) != 2 || tgIDs[0] != 111 || tgIDs[1] != 222 {
		t.Fatalf("expected [111 222], got %v", tgIDs)
	}

	// GetTelegramUserIDs panic on parse error
	assertPanic(t, "ParseUserAccount error in GetTelegramUserIDs", func() {
		badTg := AccountsOfUser{Accounts: []string{"telegram:a:b:c:d"}}
		badTg.GetTelegramUserIDs()
	})
	assertPanic(t, "ParseInt error in GetTelegramUserIDs", func() {
		badTg := AccountsOfUser{Accounts: []string{"telegram::not_a_number"}}
		badTg.GetTelegramUserIDs()
	})

	// Deprecated Get* methods returning errors
	if _, err := a.GetTelegramAccounts(); err == nil {
		t.Fatal("expected error from GetTelegramAccounts")
	}
	if _, err := a.GetGoogleAccount(); err == nil {
		t.Fatal("expected error from GetGoogleAccount")
	}
	if _, err := a.GetFbAccounts(); err == nil {
		t.Fatal("expected error from GetFbAccounts")
	}
	if _, err := a.GetFbAccount(""); err == nil {
		t.Fatal("expected error from GetFbAccount")
	}
	if _, err := a.GetFbmAccount(""); err == nil {
		t.Fatal("expected error from GetFbmAccount")
	}

	// GetAccounts
	var multi AccountsOfUser
	multi.Accounts = []string{"google:app1:u1", "google:app2:u2", "facebook:app1:u3"}
	accs, err := multi.GetAccounts("google")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(accs) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accs))
	}
	// GetAccounts parse error
	multiBad := AccountsOfUser{Accounts: []string{"google:a:b:c:d"}}
	if _, err := multiBad.GetAccounts("google"); err == nil {
		t.Fatal("expected error from GetAccounts when item cannot be parsed")
	}

	// GetAccount
	// Not found
	acc, err := multi.GetAccount("apple", "")
	if err != nil || acc != nil {
		t.Fatalf("expected nil, nil for not found, got %v, %v", acc, err)
	}
	// Found 1
	acc, err = multi.GetAccount("facebook", "app1")
	if err != nil || acc == nil || acc.ID != "u3" {
		t.Fatalf("expected u3, got %v, %v", acc, err)
	}
	// Found >1 -> returns first account and error
	multiDup := AccountsOfUser{Accounts: []string{"google::u1", "google::u2"}}
	acc, err = multiDup.GetAccount("google", "")
	if err == nil || acc == nil || acc.ID != "u1" {
		t.Fatalf("expected u1 with error for duplicate provider accounts, got %v, %v", acc, err)
	}
	// GetAccount parse error
	multiBad2 := AccountsOfUser{Accounts: []string{"google::a:b:c:d"}}
	if _, err := multiBad2.GetAccount("google", ""); err == nil {
		t.Fatal("expected error for bad account string")
	}
}

func TestWithLastLogin_Coverage(t *testing.T) {
	var w WithLastLogin
	if err := w.Validate(); err == nil {
		t.Fatal("expected error on zero LastLoginAt")
	}
	now := time.Now()
	upd := w.SetLastLoginAt(now)
	if upd.FieldName() != "lastLoginAt" {
		t.Fatalf("expected lastLoginAt field name, got %s", upd.FieldName())
	}
	if w.LastLoginAt != now {
		t.Fatalf("expected LastLoginAt to be set to %v", now)
	}
	if err := w.Validate(); err != nil {
		t.Fatalf("expected nil error on valid WithLastLogin, got %v", err)
	}
}

func TestBaseUserFields_Coverage(t *testing.T) {
	// Valid
	buf := BaseUserFields{
		NameFields: person.NameFields{
			FirstName: "Bob",
		},
	}
	if err := buf.Validate(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// Invalid NameFields
	bufBad := BaseUserFields{
		NameFields: person.NameFields{
			FirstName: " Bob",
		},
	}
	if err := bufBad.Validate(); err == nil {
		t.Fatal("expected error for invalid NameFields")
	}
}

func assertPanic(t *testing.T, msg string, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for [%s], but did not panic", msg)
		}
	}()
	f()
}
