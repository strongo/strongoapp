package strongo

import (
	"context"
	"reflect"
)

type ExecutionContext interface {
	SingleLocaleTranslator
	Context() context.Context
}

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

func NewExecutionContext(c context.Context, translator SingleLocaleTranslator) ExecutionContext {
	return executionContext{
		c: c,
		SingleLocaleTranslator: translator,
	}
}
