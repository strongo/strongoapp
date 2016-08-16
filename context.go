package strongo

import "reflect"

type ExecutionContext interface {
	SingleLocaleTranslator
	Logger() Logger
}

type AppContext interface {
	AppUserEntityKind() string
	AppUserEntityType() reflect.Type
	NewAppUserEntity() AppUser
	GetTranslator(l Logger) Translator
	SupportedLocales() LocalesProvider
}

type executionContext struct {
	SingleLocaleTranslator
	logger Logger
}

func (c executionContext) Logger() Logger {
	return c.logger
}

func NewExecutionContext(translator SingleLocaleTranslator, logger Logger) ExecutionContext {
	return executionContext{
		SingleLocaleTranslator: translator,
		logger: logger,
	}
}