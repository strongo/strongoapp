package user

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/slice"
	"strconv"
	"strings"
	"time"
)

type ownedByUser struct {
	DtCreated time.Time
	DtUpdated time.Time `datastore:",omitempty"`
}

type OwnedByUserWithIntID struct {
	ownedByUser
	AppUserIntID int64
}

func NewOwnedByUserWithIntID(id int64, created time.Time) OwnedByUserWithIntID {
	return OwnedByUserWithIntID{
		ownedByUser: ownedByUser{
			DtCreated: created,
		},
	}
}

var (
	_ BelongsToUserWithIntID = (*OwnedByUserWithIntID)(nil)
	_ CreatedTimesSetter     = (*OwnedByUserWithIntID)(nil)
	_ UpdatedTimeSetter      = (*OwnedByUserWithIntID)(nil)
)

func (ownedByUser OwnedByUserWithIntID) Validate() error {
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

func (ownedByUser *OwnedByUserWithIntID) GetAppUserID() interface{} {
	return ownedByUser.AppUserIntID
}

func (ownedByUser *OwnedByUserWithIntID) GetAppUserIntID() int64 {
	return ownedByUser.AppUserIntID
}

func (ownedByUser *OwnedByUserWithIntID) SetAppUserIntID(appUserID int64) {
	ownedByUser.AppUserIntID = appUserID
}

func (ownedByUser *ownedByUser) SetCreatedTime(v time.Time) {
	ownedByUser.DtCreated = v
}

func (ownedByUser *ownedByUser) SetUpdatedTime(v time.Time) {
	ownedByUser.DtUpdated = v
}

type AccountData interface {
	BelongsToUser
	SetLastLogin(time time.Time)
	IsEmailConfirmed() bool
	GetNames() Names
}

type AccountRecord[T comparable] interface {
	record.WithID[T]
	AccountData() AccountData
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

type Account struct {
	// Global ID of Account
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

type AccountsOfUser struct {
	// Member of TgUserEntity class
	Accounts []string `datastore:",noindex"`
}

func (ua *AccountsOfUser) AddAccount(userAccount Account) (changed bool) {
	// TODO: if !IsKnownUserAccountProvider(userAccount.Provider) {
	// 	panic("Unknown provider: " + userAccount.Provider)
	// }
	if userAccount.ID == "" || userAccount.ID == "0" {
		panic(fmt.Sprintf("Invalid userAccount.ID: [%v], userAccount.String: %v", userAccount.ID, userAccount.String()))
	} else if strings.Contains(userAccount.ID, ":") {
		panic("ID should not contains the ':' character.")
	}

	if userAccount.App == "" {
		switch userAccount.Provider {
		case "google", "email": // TODO: abstract this to provider definition
			// It's OK to have empty app for this providers
		default:
			panic(fmt.Sprintf("User account must have non-empty App field, got: %+v", userAccount))
		}
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

func (ua *AccountsOfUser) SetBotUserID(platform, botID, botUserID string) {
	ua.AddAccount(Account{
		Provider: platform,
		App:      botID,
		ID:       botUserID,
	})
}

// RemoveAccount removes account from the list of account IDs.
func (ua *AccountsOfUser) RemoveAccount(userAccount Account) (changed bool) {
	count := len(ua.Accounts)
	ua.Accounts = slice.RemoveInPlace(userAccount.String(), ua.Accounts)
	return len(ua.Accounts) != count
}

func userAccountPrefix(provider, app string) string {
	if app == "" {
		return provider + ":"
	} else {
		return provider + ":" + app + ":"
	}
}

func (ua *AccountsOfUser) HasAccount(provider, app string) bool {
	p := userAccountPrefix(provider, app)
	for _, a := range ua.Accounts {
		if strings.HasPrefix(a, p) {
			return true
		}
	}
	return false
}

func (ua *AccountsOfUser) HasTelegramAccount() bool {
	return ua.HasAccount("telegram", "")
}

func (ua *AccountsOfUser) HasGoogleAccount() bool {
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

// GetAccount returns the first account of the given provider and app.
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

// LastLogin is a struct that contains the last login time of a user.
type LastLogin struct {
	DtLastLogin time.Time `datastore:",omitempty"`
}

// SetLastLogin sets the last login time of a user.
func (l *LastLogin) SetLastLogin(time time.Time) {
	l.DtLastLogin = time
}
