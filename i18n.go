package strongo

import (
	"fmt"
	"strings"
)

// Translator is interface for translate providers
type Translator interface {
	Translate(key, locale string, args ...interface{}) string
	TranslateNoWarning(key, locale string, args ...interface{}) string
}

// SingleLocaleTranslator should be implemente by translators to a single language
type SingleLocaleTranslator interface {
	Locale() Locale
	Translate(key string, args ...interface{}) string
	TranslateNoWarning(key string, args ...interface{}) string
}

// SingleLocaleTranslatorWithBackup should be implemente by translators to a single language with backup to another one.
type SingleLocaleTranslatorWithBackup struct {
	PrimaryTranslator SingleLocaleTranslator
	BackupTranslator  SingleLocaleTranslator
}

// NewSingleLocaleTranslatorWithBackup creates SingleLocaleTranslatorWithBackup
func NewSingleLocaleTranslatorWithBackup(primary, backup SingleLocaleTranslator) SingleLocaleTranslatorWithBackup {
	return SingleLocaleTranslatorWithBackup{PrimaryTranslator: primary, BackupTranslator: backup}
}

// Locale returns local of the translator
func (t SingleLocaleTranslatorWithBackup) Locale() Locale {
	return t.PrimaryTranslator.Locale()
}

// Translate translates
func (t SingleLocaleTranslatorWithBackup) Translate(key string, args ...interface{}) string {
	result := t.PrimaryTranslator.Translate(key, args...)
	if result == key {
		result = t.BackupTranslator.Translate(key, args...)
	}
	return result
}

// TranslateNoWarning translates and does not log warning if translation not found
func (t SingleLocaleTranslatorWithBackup) TranslateNoWarning(key string, args ...interface{}) string {
	result := t.PrimaryTranslator.TranslateNoWarning(key, args...)
	if result == key {
		result = t.BackupTranslator.TranslateNoWarning(key, args...)
	}
	return result
}

// LocalesProvider provides locale by code
type LocalesProvider interface {
	GetLocaleByCode5(code5 string) (Locale, error)
}

// Locale describes language
type Locale struct {
	Code5        string
	IsRtl        bool
	NativeTitle  string
	EnglishTitle string
	FlagIcon     string
}

// SiteCode returns code for using in website URLs
func (l Locale) SiteCode() string {
	s := strings.ToLower(l.Code5)
	if s1 := s[:2]; s1 == s[3:] || s1 == "en" || s1 == "fa" || s1 == "ja" || s1 == "zh" {
		return s1
	}
	return s
}

// String represent locale information as string
func (l Locale) String() string {
	return fmt.Sprintf(`Locale{Code5: "%v", IsRtl: %v, NativeTitle: "%v", EnglishTitle: "%v", FlagIcon: "%v"}`, l.Code5, l.IsRtl, l.NativeTitle, l.EnglishTitle, l.FlagIcon)
}

// TitleWithIcon returns name of the language and flag emoji
func (l Locale) TitleWithIcon() string {
	if l.IsRtl {
		return l.NativeTitle + " " + l.FlagIcon
	}
	return l.FlagIcon + " " + l.NativeTitle
}

// TitleWithIconAndNumber returns name, flag emoji and a number // TODO: should bot be here
func (l Locale) TitleWithIconAndNumber(i int) string {
	if l.IsRtl {
		return fmt.Sprintf("%v %v .%d/", l.FlagIcon, l.NativeTitle, i)
	}
	return fmt.Sprintf("/%d. %v %v", i, l.NativeTitle, l.FlagIcon)
}
