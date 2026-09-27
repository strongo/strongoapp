package with

import (
	"testing"
	"time"
)

func TestCountryID_Coverage(t *testing.T) {
	// OptionalCountryID
	opt := OptionalCountryID{CountryID: "US"}
	if err := opt.Validate(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	optBad := OptionalCountryID{CountryID: "bad"}
	if err := optBad.Validate(); err == nil {
		t.Fatal("expected error for bad OptionalCountryID")
	}

	// RequiredCountryID
	req := RequiredCountryID{CountryID: "US"}
	if err := req.Validate(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	reqBad := RequiredCountryID{CountryID: ""}
	if err := reqBad.Validate(); err == nil {
		t.Fatal("expected error for empty RequiredCountryID")
	}

	// ValidateOptionalCountryID cases
	if err := ValidateOptionalCountryID("countryID", ""); err != nil {
		t.Fatalf("expected nil on empty, got %v", err)
	}
	if err := ValidateOptionalCountryID("countryID", " US "); err == nil {
		t.Fatal("expected error for spaces")
	}
	// Colon panic
	assertPanic(t, "colon in countryID", func() {
		_ = ValidateOptionalCountryID("countryID", "US:NY")
	})
	// Length != 2
	if err := ValidateOptionalCountryID("countryID", "USA"); err == nil {
		t.Fatal("expected error for len != 2")
	}
	// Not uppercase
	if err := ValidateOptionalCountryID("countryID", "us"); err == nil {
		t.Fatal("expected error for lowercase")
	}
	// Unknown country
	if err := ValidateOptionalCountryID("countryID", "ZZ"); err == nil {
		t.Fatal("expected error for unknown country")
	}
	// UnknownCountryID constant "--"
	if err := ValidateOptionalCountryID("countryID", UnknownCountryID); err != nil {
		t.Fatalf("expected nil for UnknownCountryID, got %v", err)
	}

	// CountryIDsField
	cids := CountryIDsField{CountryIDs: []string{"US", "GB"}}
	if err := cids.Validate(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	cidsBad := CountryIDsField{CountryIDs: []string{"US", "invalid"}}
	if err := cidsBad.Validate(); err == nil {
		t.Fatal("expected error for invalid country in slice")
	}
}

func TestCreated_Coverage(t *testing.T) {
	// Created.Validate
	c := Created{At: "2024-01-01", By: "admin"}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid date, got %v", err)
	}
	cBadDate := Created{At: "2024-99-99", By: "admin"}
	if err := cBadDate.Validate(); err == nil {
		t.Fatal("expected error for invalid DateOnly")
	}
	cRFC := Created{At: "2024-01-01T12:00:00Z", By: "admin"}
	if err := cRFC.Validate(); err != nil {
		t.Fatalf("expected valid RFC3339, got %v", err)
	}
	cBadRFC := Created{At: "invalid-rfc-timestamp", By: "admin"}
	if err := cBadRFC.Validate(); err == nil {
		t.Fatal("expected error for bad RFC3339")
	}
	cMissing := Created{At: "", By: ""}
	if err := cMissing.Validate(); err == nil {
		t.Fatal("expected error for missing at and by")
	}

	// CreatedField
	cf := CreatedField{Created: c}
	if err := cf.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
	cfBad := CreatedField{Created: cMissing}
	if err := cfBad.Validate(); err == nil {
		t.Fatal("expected error on CreatedField.Validate")
	}

	// CreatedFields
	now := time.Now()
	cfs := CreatedFields{
		CreatedAtField: CreatedAtField{CreatedAt: now},
		CreatedByField: CreatedByField{CreatedBy: "admin"},
	}
	if err := cfs.Validate(); err != nil {
		t.Fatalf("expected valid CreatedFields, got %v", err)
	}
	upds := cfs.UpdatesWhenCreated()
	if len(upds) != 2 {
		t.Fatalf("expected 2 updates, got %d", len(upds))
	}
	cfsBad := CreatedFields{}
	if err := cfsBad.Validate(); err == nil {
		t.Fatal("expected error on empty CreatedFields")
	}

	// CreatedAtField GetCreatedTime, SetCreatedAt, UpdatesCreatedOn
	cat := CreatedAtField{}
	cat.SetCreatedAt(now)
	gotTime, err := cat.GetCreatedTime()
	if err != nil || gotTime != now {
		t.Fatalf("expected %v, got %v, err: %v", now, gotTime, err)
	}
	updsCat := cat.UpdatesCreatedOn()
	if len(updsCat) != 1 || updsCat[0].FieldName() != "createdOn" {
		t.Fatalf("unexpected UpdatesCreatedOn: %v", updsCat)
	}

	// CreatedByField SetCreatedBy, GetCreatedBy, UpdatesCreatedBy
	cbf := CreatedByField{}
	cbf.SetCreatedBy("superadmin")
	if cbf.GetCreatedBy() != "superadmin" {
		t.Fatalf("expected superadmin, got %s", cbf.GetCreatedBy())
	}
	updsCbf := cbf.UpdatesCreatedBy()
	if len(updsCbf) != 1 || updsCbf[0].FieldName() != "createdBy" {
		t.Fatalf("unexpected UpdatesCreatedBy: %v", updsCbf)
	}
}

func TestDatesFields_Coverage(t *testing.T) {
	df := DatesFields{}

	// UpdatesWhenDatesChanged on empty
	updsEmpty := df.UpdatesWhenDatesChanged()
	if len(updsEmpty) != 3 {
		t.Fatalf("expected 3 updates, got %d", len(updsEmpty))
	}

	// AddDate
	upd1 := df.AddDate("2024-05-10")
	if len(upd1) != 3 { // dates, dateMax, dateMin
		t.Fatalf("expected 3 updates on first AddDate, got %d", len(upd1))
	}
	// Duplicate AddDate
	updDup := df.AddDate("2024-05-10")
	if len(updDup) != 0 {
		t.Fatalf("expected 0 updates on duplicate, got %d", len(updDup))
	}
	// Add higher date -> updates max
	updMax := df.AddDate("2024-05-20")
	if len(updMax) != 2 || df.DateMax != "2024-05-20" {
		t.Fatalf("expected dateMax update, got %v", updMax)
	}
	// Add lower date -> updates min
	updMin := df.AddDate("2024-05-01")
	if len(updMin) != 2 || df.DateMin != "2024-05-01" {
		t.Fatalf("expected dateMin update, got %v", updMin)
	}

	// UpdatesWhenDatesChanged on non-empty
	updsNonEmpty := df.UpdatesWhenDatesChanged()
	if len(updsNonEmpty) != 3 {
		t.Fatalf("expected 3 updates, got %d", len(updsNonEmpty))
	}

	// Validate valid
	if err := df.Validate(); err != nil {
		t.Fatalf("expected valid DatesFields, got %v", err)
	}
	if df.DateMin != "2024-05-01" || df.DateMax != "2024-05-20" {
		t.Fatalf("unexpected DateMin/DateMax: %s, %s", df.DateMin, df.DateMax)
	}

	// Validate empty
	dfEmpty := DatesFields{}
	if err := dfEmpty.Validate(); err != nil {
		t.Fatalf("expected valid on empty, got %v", err)
	}

	// Validate space/empty date
	dfBad1 := DatesFields{Dates: []string{""}}
	if err := dfBad1.Validate(); err == nil {
		t.Fatal("expected error for empty date item")
	}

	// Validate bad date string
	dfBad2 := DatesFields{Dates: []string{"invalid-date"}}
	if err := dfBad2.Validate(); err == nil {
		t.Fatal("expected error for invalid date format")
	}

	// Validate duplicate date
	dfBad3 := DatesFields{Dates: []string{"2024-05-10", "2024-05-10"}}
	if err := dfBad3.Validate(); err == nil {
		t.Fatal("expected error for duplicate date")
	}
}

func TestDeletedFields_Coverage(t *testing.T) {
	now := time.Now()
	df := DeletedFields{DeletedAt: now, DeletedBy: "admin"}

	// UpdatesWhenDeleted
	upds := df.UpdatesWhenDeleted()
	if len(upds) != 2 {
		t.Fatalf("expected 2 updates, got %d", len(upds))
	}

	// Validate both set
	if err := df.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}

	// Validate both zero
	dfEmpty := DeletedFields{}
	if err := dfEmpty.Validate(); err != nil {
		t.Fatalf("expected valid on empty, got %v", err)
	}

	// Validate DeletedAt set, DeletedBy empty
	dfMissingBy := DeletedFields{DeletedAt: now}
	if err := dfMissingBy.Validate(); err == nil {
		t.Fatal("expected error on missing DeletedBy")
	}

	// Validate DeletedBy set, DeletedAt zero
	dfMissingAt := DeletedFields{DeletedBy: "admin"}
	if err := dfMissingAt.Validate(); err == nil {
		t.Fatal("expected error on missing DeletedAt")
	}
}

func TestEmailsField_Coverage(t *testing.T) {
	ef := EmailsField{
		Emails: map[string]*CommunicationChannelProps{
			"valid@example.com": {
				CreatedFields: CreatedFields{
					CreatedAtField: CreatedAtField{CreatedAt: time.Now()},
					CreatedByField: CreatedByField{CreatedBy: "admin"},
				},
			},
		},
	}
	if err := ef.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}

	efBad := EmailsField{
		Emails: map[string]*CommunicationChannelProps{
			"not-an-email": {
				CreatedFields: CreatedFields{
					CreatedAtField: CreatedAtField{CreatedAt: time.Now()},
					CreatedByField: CreatedByField{CreatedBy: "admin"},
				},
			},
		},
	}
	if err := efBad.Validate(); err == nil {
		t.Fatal("expected error for invalid email address")
	}
}

func TestFlagsField_Coverage(t *testing.T) {
	ff := FlagsField{Flags: []string{"flag1", "flag2"}}
	if err := ff.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
	if ff.String() != "flags=flag1,flag2" {
		t.Fatalf("expected flags=flag1,flag2, got %s", ff.String())
	}

	ffBad := FlagsField{Flags: []string{"flag1", "  "}}
	if err := ffBad.Validate(); err == nil {
		t.Fatal("expected error for empty flag")
	}
}

func TestKeysField_Coverage(t *testing.T) {
	// UpdatesWhenKeysChanged
	kfEmpty := KeysField{}
	updsEmpty := kfEmpty.UpdatesWhenKeysChanged()
	if len(updsEmpty) != 1 {
		t.Fatalf("expected 1 update, got %d", len(updsEmpty))
	}

	kf := KeysField{Keys: []string{"k1", "k2"}}
	upds := kf.UpdatesWhenKeysChanged()
	if len(upds) != 1 {
		t.Fatalf("expected 1 update, got %d", len(upds))
	}

	// Validate valid
	if err := kf.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}

	// Validate empty key
	kfBad1 := KeysField{Keys: []string{""}}
	if err := kfBad1.Validate(); err == nil {
		t.Fatal("expected error for empty key")
	}

	// Validate untrimmed key
	kfBad2 := KeysField{Keys: []string{" k1 "}}
	if err := kfBad2.Validate(); err == nil {
		t.Fatal("expected error for untrimmed key")
	}

	// Validate duplicate key
	kfBad3 := KeysField{Keys: []string{"k1", "k1"}}
	if err := kfBad3.Validate(); err == nil {
		t.Fatal("expected error for duplicate key")
	}
}

func TestPreferredLocaleField_Coverage(t *testing.T) {
	pl := &PreferredLocaleField{}
	if err := pl.SetPreferredLocale("en-US"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pl.GetPreferredLocale() != "en-US" {
		t.Fatalf("expected en-US, got %s", pl.GetPreferredLocale())
	}
}

func TestRolesField_Coverage(t *testing.T) {
	rf := &RolesField{Roles: []string{"admin", "editor"}}

	// HasRole
	if !rf.HasRole("admin") || rf.HasRole("viewer") {
		t.Fatal("unexpected HasRole results")
	}

	// AddRole
	if rf.AddRole("admin") {
		t.Fatal("expected false when adding existing role")
	}
	if !rf.AddRole("viewer") {
		t.Fatal("expected true when adding new role")
	}

	// RemoveRole
	upds := rf.RemoveRole("viewer")
	if len(upds) != 1 {
		t.Fatalf("expected 1 update on removal, got %d", len(upds))
	}
	updsNotFound := rf.RemoveRole("not_found")
	if len(updsNotFound) != 0 {
		t.Fatalf("expected 0 updates when not found, got %d", len(updsNotFound))
	}

	// Validate
	if err := rf.Validate(); err != nil {
		t.Fatalf("expected valid RolesField, got %v", err)
	}
	rfBad := &RolesField{Roles: []string{"admin", "admin"}}
	if err := rfBad.Validate(); err == nil {
		t.Fatal("expected error for duplicate role in Validate")
	}
}

func TestTagsField_Coverage(t *testing.T) {
	tf := TagsField{Tags: []string{"tag1", "tag2"}}
	if err := tf.Validate(); err != nil {
		t.Fatalf("expected valid TagsField, got %v", err)
	}
	if tf.String() != "tags=tag1,tag2" {
		t.Fatalf("expected tags=tag1,tag2, got %s", tf.String())
	}

	tfBad := TagsField{Tags: []string{"tag1", "tag1"}}
	if err := tfBad.Validate(); err == nil {
		t.Fatal("expected error on duplicate tags")
	}
}

func TestUpdatedFields_Coverage(t *testing.T) {
	now := time.Now()
	uf := &UpdatedFields{}

	// SetUpdatedTime & GetUpdatedTime
	if err := uf.SetUpdatedTime(now); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uf.GetUpdatedTime() != now {
		t.Fatalf("expected %v, got %v", now, uf.GetUpdatedTime())
	}

	// UpdatesWhenUpdatedFieldsChanged
	uf.UpdatedBy = "user1"
	upds := uf.UpdatesWhenUpdatedFieldsChanged()
	if len(upds) != 2 {
		t.Fatalf("expected 2 updates, got %d", len(upds))
	}

	// Validate valid
	if err := uf.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}

	// Validate missing UpdatedAt
	ufBad1 := &UpdatedFields{UpdatedBy: "user1"}
	if err := ufBad1.Validate(); err == nil {
		t.Fatal("expected error for missing UpdatedAt")
	}

	// Validate missing UpdatedBy
	ufBad2 := &UpdatedFields{UpdatedAt: now}
	if err := ufBad2.Validate(); err == nil {
		t.Fatal("expected error for missing UpdatedBy")
	}
}

func TestProps_And_Validation_Branches(t *testing.T) {
	now := time.Now()
	validProps := &CommunicationChannelProps{
		CreatedFields: CreatedFields{
			CreatedAtField: CreatedAtField{CreatedAt: now},
			CreatedByField: CreatedByField{CreatedBy: "admin"},
		},
	}

	// 1. Key with leading/trailing spaces in validateCommunicationChannelsField
	efSpaceKey := EmailsField{
		Emails: map[string]*CommunicationChannelProps{
			" test@example.com ": validProps,
		},
	}
	if err := efSpaceKey.Validate(); err == nil {
		t.Fatal("expected error for space key")
	}

	// 2. p.Validate() error in validateCommunicationChannelsField
	badProps := &CommunicationChannelProps{}
	efBadProps := EmailsField{
		Emails: map[string]*CommunicationChannelProps{
			"test@example.com": badProps,
		},
	}
	if err := efBadProps.Validate(); err == nil {
		t.Fatal("expected error for bad props")
	}

	// 3. Multiple primary emails
	p1 := *validProps
	p1.IsPrimary = true
	p2 := *validProps
	p2.IsPrimary = true
	efTwoPrimary := EmailsField{
		Emails: map[string]*CommunicationChannelProps{
			"a@example.com": &p1,
			"b@example.com": &p2,
		},
	}
	if err := efTwoPrimary.Validate(); err == nil {
		t.Fatal("expected error for multiple primary emails")
	}

	// 4. p.Original == k
	pSameOrig := *validProps
	pSameOrig.Original = "same@example.com"
	efSameOrig := EmailsField{
		Emails: map[string]*CommunicationChannelProps{
			"same@example.com": &pSameOrig,
		},
	}
	if err := efSameOrig.Validate(); err == nil {
		t.Fatal("expected error for original same as key")
	}

	// 5. CommunicationChannelProps.Validate branches:
	// 5a. CreatedFields error
	pCreatedErr := CommunicationChannelProps{}
	if err := pCreatedErr.Validate(); err == nil {
		t.Fatal("expected CreatedFields error")
	}

	// 5b. TagsField error
	pTagsErr := *validProps
	pTagsErr.TagsField = TagsField{Tags: []string{"tag1", "tag1"}}
	if err := pTagsErr.Validate(); err == nil {
		t.Fatal("expected TagsField error")
	}

	// 5c. Original has spaces but empty
	pOrigSpace := *validProps
	pOrigSpace.Original = "   "
	if err := pOrigSpace.Validate(); err == nil {
		t.Fatal("expected error for original spaces but empty")
	}

	// 5d. Original has leading/trailing spaces
	pOrigTrim := *validProps
	pOrigTrim.Original = " email@example.com "
	if err := pOrigTrim.Validate(); err == nil {
		t.Fatal("expected error for original with trailing spaces")
	}

	// 5e. Type has spaces
	pTypeSpace := *validProps
	pTypeSpace.Type = " personal "
	if err := pTypeSpace.Validate(); err == nil {
		t.Fatal("expected error for type with spaces")
	}

	// 5f. Title has spaces
	pTitleSpace := *validProps
	pTitleSpace.Title = " title "
	if err := pTitleSpace.Validate(); err == nil {
		t.Fatal("expected error for title with spaces")
	}

	// 5g. Note has spaces
	pNoteSpace := *validProps
	pNoteSpace.Note = " note "
	if err := pNoteSpace.Validate(); err == nil {
		t.Fatal("expected error for note with spaces")
	}

	// 6. ValidateDateString with 10 chars but invalid date
	if _, err := ValidateDateString("2024-99-99"); err == nil {
		t.Fatal("expected error for invalid date")
	}

	// 7. ValidateRecordID branches
	if err := ValidateRecordID(""); err == nil {
		t.Fatal("expected error for empty id")
	}
	if err := ValidateRecordID(" id "); err == nil {
		t.Fatal("expected error for spaces around id")
	}
	if err := ValidateRecordID("id with space"); err == nil {
		t.Fatal("expected error for internal space")
	}
	if err := ValidateRecordID("valid-id_123"); err != nil {
		t.Fatalf("expected valid id, got %v", err)
	}

	// 8. ValidateSetSliceField branches
	if err := ValidateSetSliceField("field", []string{""}, false); err == nil {
		t.Fatal("expected error for empty slice item")
	}
	if err := ValidateSetSliceField("field", []string{" val "}, false); err == nil {
		t.Fatal("expected error for untrimmed slice item")
	}
	if err := ValidateSetSliceField("field", []string{"bad id"}, true); err == nil {
		t.Fatal("expected error for bad record id in slice")
	}

	// 9. CommChannelFields GetCommChannels lazy init and default
	ccf := &CommChannelFields{}
	channels, name := ccf.GetCommChannels(CommChannelTypeEmail)
	if channels == nil || name != EmailsFieldName {
		t.Fatalf("expected initialized emails, got %v, %s", channels, name)
	}
	channelsPhone, namePhone := ccf.GetCommChannels(CommChannelTypePhone)
	if channelsPhone == nil || namePhone != PhonesFieldName {
		t.Fatalf("expected initialized phones, got %v, %s", channelsPhone, namePhone)
	}
	unknownChannels, unknownName := ccf.GetCommChannels("unknown")
	if unknownChannels != nil || unknownName != "" {
		t.Fatalf("expected nil for unknown, got %v, %s", unknownChannels, unknownName)
	}

	// 10. Created.Validate with valid date but empty by (line 35)
	cValidAtEmptyBy := Created{At: "2024-01-01", By: ""}
	if err := cValidAtEmptyBy.Validate(); err == nil {
		t.Fatal("expected error when At is valid but By is empty")
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

