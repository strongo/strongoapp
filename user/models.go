package user

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/strongo/app/slices"
	"github.com/strongo/db"
	"strconv"
	"strings"
	"time"
)

type OwnedByUser struct { // TODO: Move to strongo/app?
	AppUserIntID int64
	DtCreated    time.Time
	DtUpdated    time.Time `datastore:",omitempty"`
}

func (ownedByUser OwnedByUser) Validate() error {
	if ownedByUser.AppUserIntID == 0 {
		return errors.New("AppUserIntID == 0")
	}
	if ownedByUser.DtCreated.IsZero() {
		return errors.New("DtCreated.IsZero()")
	}
	if ownedByUser.DtUpdated.IsZero() {
		ownedByUser.DtUpdated = ownedByUser.DtCreated
	} else if ownedByUser.DtUpdated.Before(ownedByUser.DtCreated) {
		return errors.New("DtUpdated.Before(DtCreated) is true")
	}
	return nil
}

func (ownedByUser *OwnedByUser) GetAppUserIntID() int64 {
	return ownedByUser.AppUserIntID
}

func (ownedByUser *OwnedByUser) SetAppUserIntID(appUserID int64) {
	ownedByUser.AppUserIntID = appUserID
}

func (ownedByUser *OwnedByUser) SetDtCreated(v time.Time) {
	ownedByUser.DtCreated = v
}

func (ownedByUser *OwnedByUser) SetDtUpdated(v time.Time) {
	ownedByUser.DtUpdated = v
}

type AccountEntity interface {
	BelongsToUser
	SetLastLogin(time time.Time)
	IsEmailConfirmed() bool
	GetNames() Names
}

type AccountRecord interface {
	AccountEntity
	db.EntityHolder
	UserAccount() Account
	GetEmail() string
}

type Names struct {
	FirstName string `datastore:",noindex"`
	LastName  string `datastore:",noindex"`
	NickName  string `datastore:",noindex"`
}

func (entity Names) GetNames() Names {
	return entity
}

type Account struct { // Global ID of Account
	Provider string
	App      string
	ID       string
}

func ParseUserAccount(s string) (ua Account, err error) {
	vals := strings.Split(s, ":")
	switch len(vals) {
	case 3:
		ua = Account{
			Provider: vals[0],
			App:      vals[1],
			ID:       vals[2],
		}
	case 2:
		ua = Account{
			Provider: vals[0],
			App:      "",
			ID:       vals[1],
		}
	default:
		err = fmt.Errorf("invalid Account string, expected 1 or 2 ':' characters, got: %d", strings.Count(s, ":"))
	}
	return
}

func (ua Account) String() string {
	return ua.Provider + ":" + ua.App + ":" + ua.ID
}

type AccountsOfUser struct {// Member of TgUserEntity class
	Accounts []string `datastore:",noindex"`
}

func (ua *AccountsOfUser) AddAccount(userAccount Account) (changed bool) {
	// TODO: if !IsKnownUserAccountProvider(userAccount.Provider) {
	//	panic("Unknown provider: " + userAccount.Provider)
	//}
	if userAccount.ID == "" || userAccount.ID == "0" {
		panic(fmt.Sprintf("Invalid userAccount.ID: [%v], userAccount.String: %v", userAccount.ID, userAccount.String()))
	} else if strings.Contains(userAccount.ID, ":") {
		panic("ID should not contains the ':' character.")
	}

	if userAccount.App == "" {
		panic(fmt.Sprintf("User account must have non-empty App field, got: %+v", userAccount))
	} else if strings.Contains(userAccount.App, ":") {
		panic("App name should not contains the ':' character.")
	}

	if userAccount.Provider == "" {
		panic("User account must have non-empty Provider field")
	} else if strings.Contains(userAccount.Provider, ":") {
		panic("Provider should not contains the ':' character.")
	}

	account := userAccount.String()
	for _, a := range ua.Accounts {
		if a == account {
			return false
		}
	}
	ua.Accounts = append(ua.Accounts, account)
	return true
}

func (ua *AccountsOfUser) RemoveAccount(userAccount Account) (changed bool) {
	ua.Accounts, changed = slices.FilterStrings(ua.Accounts, []string{userAccount.String()})
	return true
}

func userAccountPrefix(provider, app string) string {
	if app == "" {
		return provider + ":"
	} else {
		return provider + ":" + app + ":"
	}
}

func (ua AccountsOfUser) HasAccount(provider, app string) bool {
	p := userAccountPrefix(provider, app)
	for _, a := range ua.Accounts {
		if strings.HasPrefix(a, p) {
			return true
		}
	}
	return false
}

func (ua AccountsOfUser) HasTelegramAccount() bool {
	return ua.HasAccount("telegram", "")
}

func (ua AccountsOfUser) HasGoogleAccount() bool {
	return ua.HasAccount("google", "")
}

func (ua *AccountsOfUser) GetTelegramUserIDs() (telegramUserIDs []int64) {
	for _, a := range ua.Accounts {
		if strings.HasPrefix(a, "telegram:") {
			if ua, err := ParseUserAccount(a); err != nil {
				panic(err)
			} else if telegramUserID, err := strconv.ParseInt(ua.ID, 10, 64); err != nil {
				panic(err)
			} else {
				telegramUserIDs = append(telegramUserIDs, telegramUserID)
			}
		}
	}
	return
}

func (ua *AccountsOfUser) GetTelegramAccounts() (telegramAccounts []Account) {
	for _, a := range ua.Accounts {
		if strings.HasPrefix(a, "telegram:") {
			if ua, err := ParseUserAccount(a); err != nil {
				panic(err)
			} else {
				telegramAccounts = append(telegramAccounts, ua)
			}
		}
	}
	return
}

func (ua *AccountsOfUser) GetGoogleAccount() (userAccount *Account, err error) {
	return ua.GetAccount("google", "")
}

func (ua *AccountsOfUser) GetFbAccounts() (userAccounts []Account, err error) {
	return ua.GetAccounts("fb")
}

func (ua *AccountsOfUser) GetAccounts(platform string) (userAccounts []Account, err error) {
	var userAccount Account
	prefix := platform + ":"
	for _, a := range ua.Accounts {
		if strings.HasPrefix(a, prefix) {
			if userAccount, err = ParseUserAccount(a); err != nil {
				return
			}
			userAccounts = append(userAccounts, userAccount)
		}
	}
	return
}

func (ua *AccountsOfUser) GetFbAccount(app string) (userAccount *Account, err error) {
	if app == "" {
		return nil, errors.New("Parameter app is required")
	}
	return ua.GetAccount("fb", app)
}

func (ua *AccountsOfUser) GetFbmAccount(fbPageID string) (userAccount *Account, err error) {
	return ua.GetAccount("fbm", fbPageID)
}

func (ua *AccountsOfUser) GetAccount(provider, app string) (userAccount *Account, err error) {
	count := 0
	prefix := userAccountPrefix(provider, app)

	for _, a := range ua.Accounts {
		if strings.HasPrefix(a, prefix) {
			if count == 0 {
				var ua Account
				if ua, err = ParseUserAccount(a); err != nil {
					return
				}
				userAccount = &ua
			}
			count += 1
		}
	}
	if userAccount != nil && count > 1 {
		err = fmt.Errorf("User has %d linked '%v' accounts", count, provider)
	}
	return
}

type LastLogin struct {
	DtLastLogin time.Time
}

func (l *LastLogin) SetLastLogin(time time.Time) {
	l.DtLastLogin = time
}
