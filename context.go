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
