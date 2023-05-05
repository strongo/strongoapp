package strongo

import (
	"context"
	"github.com/strongo/i18n"
	"reflect"
)

// ExecutionContext is execution context for UI app
// Contains translator and context.Context
type ExecutionContext interface {
	i18n.SingleLocaleTranslator
	Context() context.Context
}

// AppContext is app user context for an app
type AppContext interface {
	AppUserEntityKind() string
	AppUserEntityType() reflect.Type
	NewAppUserEntity() AppUser
	GetTranslator(c context.Context) i18n.Translator
	SupportedLocales() i18n.LocalesProvider
}

type executionContext struct {
	c context.Context
	i18n.SingleLocaleTranslator
}

func (ec executionContext) Context() context.Context {
	return ec.c
}

// NewExecutionContext creates new execution context
func NewExecutionContext(c context.Context, translator i18n.SingleLocaleTranslator) ExecutionContext {
	return executionContext{
		c:                      c,
		SingleLocaleTranslator: translator,
	}
}
