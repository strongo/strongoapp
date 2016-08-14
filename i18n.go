package strongo

import (
	"fmt"
	"strings"
)

type Translator interface {
	Translate(key, locale string, args ...interface{}) string
	TranslateNoWarning(key, locale string, args ...interface{}) string
}

type SingleLocaleTranslator interface {
	Locale() Locale
	Translate(key string, args ...interface{}) string
	TranslateNoWarning(key string, args ...interface{}) string
}

type SingleLocaleTranslatorWithBackup struct {
	PrimaryTranslator SingleLocaleTranslator
	BackupTranslator SingleLocaleTranslator
}

func NewSingleLocaleTranslatorWithBackup(primary, backup SingleLocaleTranslator) SingleLocaleTranslatorWithBackup {
	return SingleLocaleTranslatorWithBackup{PrimaryTranslator: primary, BackupTranslator: backup}
}

func (t SingleLocaleTranslatorWithBackup) Locale() Locale {
	return t.PrimaryTranslator.Locale()
}

func (t SingleLocaleTranslatorWithBackup) Translate(key string, args ...interface{}) string {
	result := t.PrimaryTranslator.Translate(key, args...)
	if result == key {
		result = t.BackupTranslator.Translate(key, args...)
	}
	return result
}

func (t SingleLocaleTranslatorWithBackup) TranslateNoWarning(key string, args ...interface{}) string {
	result := t.PrimaryTranslator.TranslateNoWarning(key, args...)
	if result == key {
		result = t.BackupTranslator.TranslateNoWarning(key, args...)
	}
	return result
}


type LocalesProvider interface {
	GetLocaleByCode5(code5 string) (Locale, error)
}

type Locale struct {
	Code5        string
	IsRtl        bool
	NativeTitle  string
	EnglishTitle string
	FlagIcon     string
}

func (l Locale) SiteCode() string {
	s := strings.ToLower(l.Code5)
	if s1 := s[:2]; s1 == s[3:] || s1 == "en" || s1 == "fa" || s1 == "ja" || s1 == "zh" {
		return s1
	}
	return s
}

func (l Locale) String() string {
	return fmt.Sprintf(`Locale{Code5: "%v", IsRtl: %v, NativeTitle: "%v", EnglishTitle: "%v", FlagIcon: "%v"}`, l.Code5, l.IsRtl, l.NativeTitle, l.EnglishTitle, l.FlagIcon)
}

func (l Locale) TitleWithIcon() string {
	if l.IsRtl {
		return l.NativeTitle + " " + l.FlagIcon
	} else {
		return l.FlagIcon + " " + l.NativeTitle
	}

}

func (l Locale) TitleWithIconAndNumber(i int) string {
	if l.IsRtl {
		return fmt.Sprintf("%v %v .%d/", l.FlagIcon, l.NativeTitle, i)
	} else {
		return fmt.Sprintf("/%d. %v %v", i, l.NativeTitle, l.FlagIcon)
	}
}
