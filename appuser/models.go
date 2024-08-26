package appuser

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"strconv"
	"strings"
	"time"
)

// NewOwnedByUserWithID creates a new OwnedByUserWithID, takes user ID and time of creation
func NewOwnedByUserWithID(id string, created time.Time) OwnedByUserWithID {
	if id == "" {
		panic("id is empty string")
	}
	if created.IsZero() {
		panic("created is zero time")
	}
	result := OwnedByUserWithID{
		AppUserID: id,
	}
	result.SetCreatedAt(created)
	return result
}

// OwnedByUserWithID is a struct that implements BelongsToUser & BelongsToUserWithIntID
type OwnedByUserWithID struct {
	AppUserID string // intentionally indexed & do NOT omitempty (so we can find records with empty AppUserID)

	// AppUserIntID is a strongly typed integer ID of a user
	// Deprecated: use AppUserID instead. Remove BelongsToUserWithIntID once AppUserIntID is removed.
	//AppUserIntID int64

	with.CreatedFields
	with.UpdatedFields
}

func (ownedByUser *OwnedByUserWithID) Validate() error {
	//if ownedByUser.AppUserIntID == 0 {
	//	return errors.New("AppUserIntID == 0")
	//}
	if ownedByUser.AppUserID == "" {
		return errors.New("AppUserID is required field")
	}
	if ownedByUser.CreatedAt.IsZero() {
		return errors.New("DtCreated.IsZero()")
	}
	if ownedByUser.UpdatedAt.IsZero() {
		ownedByUser.UpdatedAt = ownedByUser.CreatedAt
	} else if ownedByUser.UpdatedAt.Before(ownedByUser.CreatedAt) {
		return errors.New("DtUpdated.Before(DtCreated) is true")
	}
	return nil
}

func (ownedByUser *OwnedByUserWithID) GetAppUserID() string {
	return ownedByUser.AppUserID
	//if ownedByUser.AppUserID != "" {
	//	return ownedByUser.AppUserID
	//}
	//return strconv.FormatInt(ownedByUser.AppUserIntID, 10)
}

//func (ownedByUser *OwnedByUserWithID) GetAppUserIntID() int64 {
//	return ownedByUser.AppUserIntID
//}

func (ownedByUser *OwnedByUserWithID) SetAppUserIntID(appUserID int64) {
	ownedByUser.SetAppUserID(strconv.FormatInt(appUserID, 10))
}

func (ownedByUser *OwnedByUserWithID) SetAppUserID(appUserID string) {
	ownedByUser.AppUserID = appUserID
	//ownedByUser.AppUserIntID = 0
}

var _ AccountData = (*AccountDataBase)(nil)

func NewEmailData(email string) EmailData {
	email = strings.TrimSpace(email)
	return EmailData{
		EmailRaw:       email,
		EmailLowerCase: strings.ToLower(email),
	}
}

// EmailData stores info about email
type EmailData struct {
	EmailRaw       string `firestore:"emailRaw"`
	EmailLowerCase string `firestore:"emailLowerCase"`
	EmailConfirmed bool   `firestore:"emailConfirmed"`
}

func (ed *EmailData) Validate() error {
	if strings.ToLower(ed.EmailRaw) == ed.EmailLowerCase {
		ed.EmailRaw = ""
	}
	if strings.ToLower(ed.EmailRaw) != ed.EmailLowerCase {
		return fmt.Errorf("EmailRaw != EmailLowerCase: %s != %s", ed.EmailRaw, ed.EmailLowerCase)
	}
	return nil
}

func (ed *EmailData) GetEmailRaw() string {
	if ed.EmailRaw != "" {
		return ed.EmailRaw
	}
	if ed.EmailLowerCase != "" {
		return ed.EmailLowerCase
	}
	return ""
}

func (ed *EmailData) GetEmailLowerCase() string {
	return ed.EmailLowerCase
}

func (ed *EmailData) GetEmailConfirmed() bool {
	return ed.EmailConfirmed
}

func (ed *EmailData) SetEmailConfirmed(value bool) {
	ed.EmailConfirmed = value
}

type AccountDataBase struct {
	AccountKey
	OwnedByUserWithID
	person.NameFields
	WithLastLogin
	EmailData

	Domains []string `json:"domains" dalgo:"domains" firestore:"domains"` // E.g. website domain names used to authenticate user

	Admin bool

	// ClientID is an OAuth2 client ID
	ClientID string `json:"clientID" dalgo:"clientID" firestore:"clientID"`

	FederatedIdentity string `firestore:"federatedIdentity"`
	FederatedProvider string `firestore:"federatedProvider"`
}

func (v *AccountDataBase) Validate() error {
	if err := v.OwnedByUserWithID.Validate(); err != nil {
		return err
	}
	if err := v.NameFields.Validate(); err != nil {
		return err
	}
	if err := v.WithLastLogin.Validate(); err != nil {
		return err
	}
	if err := v.EmailData.Validate(); err != nil {
		return err
	}
	return nil
}

func (v *AccountDataBase) GetNames() person.NameFields {
	//TODO implement me
	panic("implement me")
}

// AccountData stores info about a user account with auth provider
type AccountData interface {
	BelongsToUser
	GetEmailLowerCase() string
	GetEmailConfirmed() bool
	SetLastLoginAt(time time.Time) dal.Update
	GetNames() person.NameFields
}

type AccountRecord interface {
	AccountKey() AccountKey
	AccountData() AccountData
}

// AccountKey stores info about a user account with auth provider
type AccountKey struct {
	// Global ID of AccountKey
	Provider string `json:"provider" dalgo:"provider" firestore:"dalgo"` // E.g. Email, Google, Facebook, etc.
	App      string `json:"app" dalgo:"app" firestore:"app"`             // E.g. app ID, bot ID, etc.
	ID       string `json:"id" dalgo:"id" firestore:"id"`                // An ID of a user at auth provider. E.g. email address, some ID, etc.
}

func (ua AccountKey) String() string {
	return ua.Provider + ":" + ua.App + ":" + ua.ID
}

func ParseUserAccount(s string) (ua AccountKey, err error) {
	vals := strings.Split(s, ":")
	switch len(vals) {
	case 3:
		ua = AccountKey{
			Provider: vals[0],
			App:      vals[1],
			ID:       vals[2],
		}
	case 2:
		ua = AccountKey{
			Provider: vals[0],
			App:      "",
			ID:       vals[1],
		}
	default:
		err = fmt.Errorf("invalid AccountKey string, expected 1 or 2 ':' characters, got: %d", strings.Count(s, ":"))
	}
	return
}

type AccountsOfUser struct {
	Accounts []string `datastore:",noindex"`
}

func (ua *AccountsOfUser) AddAccount(userAccount AccountKey) (changed bool) {
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
	ua.AddAccount(AccountKey{
		Provider: platform,
		App:      botID,
		ID:       botUserID,
	})
}

// RemoveAccount removes an account from the list of account IDs.
func (ua *AccountsOfUser) RemoveAccount(userAccount AccountKey) (changed bool) {
	count := len(ua.Accounts)
	ua.Accounts = removeInPlace(userAccount.String(), ua.Accounts)
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

// Deprecated: use HasAccount instead
func (ua *AccountsOfUser) HasTelegramAccount() bool {
	panic("deprecated")
}

// Deprecated: use HasAccount instead
func (ua *AccountsOfUser) HasGoogleAccount() bool {
	panic("deprecated")
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

// Deprecated: use GetAccounts instead
func (ua *AccountsOfUser) GetTelegramAccounts() (telegramAccounts []AccountKey, er error) {
	return nil, errors.New("GetTelegramAccounts() is deprecated, use GetAccounts(platform string) instead")
}

// Deprecated: use GetAccounts instead
func (ua *AccountsOfUser) GetGoogleAccount() (userAccount *AccountKey, err error) {
	return nil, errors.New("GetGoogleAccount() is deprecated, use GetAccount(provider, app string) instead")
}

// Deprecated: use GetAccounts instead
func (ua *AccountsOfUser) GetFbAccounts() (userAccounts []AccountKey, err error) {
	return nil, errors.New("GetFbAccounts() is deprecated, use GetAccounts(platform string) instead")
}

func (ua *AccountsOfUser) GetAccounts(platform string) (userAccounts []AccountKey, err error) {
	var userAccount AccountKey
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

// Deprecated: use GetAccount instead
func (ua *AccountsOfUser) GetFbAccount(app string) (userAccount *AccountKey, err error) {
	return nil, errors.New("GetFbAccount() is deprecated, use GetAccount() instead")
}

// Deprecated: use GetAccount instead
func (ua *AccountsOfUser) GetFbmAccount(fbPageID string) (userAccount *AccountKey, err error) {
	return nil, errors.New("GetFbmAccount() is deprecated, use GetAccount() instead")
}

// GetAccount returns the first account of the given provider and app.
func (ua *AccountsOfUser) GetAccount(provider, app string) (userAccount *AccountKey, err error) {
	count := 0
	prefix := userAccountPrefix(provider, app)

	for _, a := range ua.Accounts {
		if strings.HasPrefix(a, prefix) {
			if count == 0 {
				var ua AccountKey
				if ua, err = ParseUserAccount(a); err != nil {
					return
				}
				userAccount = &ua
			}
			count += 1
		}
	}
	if userAccount != nil && count > 1 {
		err = fmt.Errorf("only 1 account from same auth provider allowed per user and user linked %d '%v' accounts", count, provider)
	}
	return
}
