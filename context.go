package strongo

import (
	"context"
	"reflect"
)

// ExecutionContext is execution context for UI app
// Contains translator and context.Context
type ExecutionContext interface {
	SingleLocaleTranslator
	Context() context.Context
}

// AppContext is app user context for an app
type AppContext interface {
	AppUserEntityKind() string
	AppUserEntityType() reflect.Type
	NewAppUserEntity() AppUser
	GetTranslator(c context.Context) Translator
	SupportedLocales() LocalesProvider
}

type executionContext struct {
	c context.Context
	SingleLocaleTranslator
}

func (ec executionContext) Context() context.Context {
	return ec.c
}

// NewExecutionContext creates new execution context
func NewExecutionContext(c context.Context, translator SingleLocaleTranslator) ExecutionContext {
	return executionContext{
		c: c,
		SingleLocaleTranslator: translator,
	}
}
