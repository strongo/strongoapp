package strongo

import (
	"bytes"
	"fmt"
	"golang.org/x/net/context"
	"html/template"
	"strings"
	"github.com/strongo/log"
)

type mapTranslator struct {
	c                 context.Context
	defaultLocale     string
	translations      map[string]map[string]string
	templatesByLocale map[string]*template.Template
}

func NewMapTranslator(c context.Context, translations map[string]map[string]string) Translator {
	return mapTranslator{
		c: c,
		defaultLocale: "en-US",
		translations: translations,
		templatesByLocale: make(map[string]*template.Template),
	}
}

type theSingleLocaleTranslator struct {
	locale Locale
	Translator
}

func (t theSingleLocaleTranslator) Translate(key string, args ...interface{}) string {
	return t.Translator.Translate(key, t.locale.Code5, args...)
}

func (t theSingleLocaleTranslator) Locale() Locale {
	return t.locale
}

func (t theSingleLocaleTranslator) TranslateNoWarning(key string, args ...interface{}) string {
	return t.Translator.TranslateNoWarning(key, t.locale.Code5, args...)
}

var _ SingleLocaleTranslator = (*theSingleLocaleTranslator)(nil)

func NewSingleMapTranslator(locale Locale, translator Translator) SingleLocaleTranslator {
	return theSingleLocaleTranslator{
		locale:     locale,
		Translator: translator,
	}
}

func (t mapTranslator) _translate(warn bool, key, locale string, args ...interface{}) string {
	s, found := t.translations[key][locale]
	if !found {
		if warn {
			log.Warningf(t.c, "Translation not found by key & locale: key=%v&locale=%v", key, locale)
		}
		if t.defaultLocale != "" {
			if s, found = t.translations[key]["en-US"]; !found {
				if warn {
					log.Errorf(t.c, "Translation not found for default locale: key=%v&locale=%v", key, locale)
				}
				return key
			}
		}
	}
	if len(args) > 0 {
		if len(args) == 1 && strings.Contains(s, "}}") && (strings.Contains(s, "{{.") || strings.Contains(s, "{{ .")) {
			tk := locale + key
			tmpl, ok := t.templatesByLocale[tk]
			if !ok {
				var err error
				if tmpl, err = template.New(key).Parse(s); err != nil {
					panic(fmt.Sprintf("Failed to parse template '%v' for locale '%v': %v", key, locale, err.Error()))
				}
				t.templatesByLocale[tk] = tmpl
			}
			var buffer bytes.Buffer
			if err := tmpl.Execute(&buffer, args[0]); err != nil {
				panic(fmt.Sprintf("Failed to render template '%v' for locale '%v': %v", key, locale, err.Error()))
			} else {
				return buffer.String()
			}
		}
		s = fmt.Sprintf(s, args...)
	}
	return s
}

func (t mapTranslator) Translate(key, locale string, args ...interface{}) string {
	return t._translate(true, key, locale, args...)
}

func (t mapTranslator) TranslateNoWarning(key, locale string, args ...interface{}) string {
	return t._translate(false, key, locale, args...)
}
