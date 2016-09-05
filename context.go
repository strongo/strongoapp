package strongo

import (
	"golang.org/x/net/context"
	"reflect"
)

type ExecutionContext interface {
	SingleLocaleTranslator
	Logger() Logger
	Context() context.Context
}

type AppContext interface {
	AppUserEntityKind() string
	AppUserEntityType() reflect.Type
	NewAppUserEntity() AppUser
	GetTranslator(l Logger) Translator
	SupportedLocales() LocalesProvider
}

type executionContext struct {
	c context.Context
	SingleLocaleTranslator
	logger Logger
}

func (ec executionContext) Logger() Logger {
	return ec.logger
}

func (ec executionContext) Context() context.Context {
	return ec.c
}

func NewExecutionContext(c context.Context, translator SingleLocaleTranslator, logger Logger) ExecutionContext {
	return executionContext{
		c: c,
		SingleLocaleTranslator: translator,
		logger:                 logger,
	}
}
