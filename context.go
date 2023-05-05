package strongo

import (
	"context"
)

// ExecutionContext is execution context for UI app
// Contains translator and context.Context
type ExecutionContext interface {
	Context() context.Context
}

// AppContext is app context for an app
// Deprecated: use AppUserSettings instead
type AppContext interface {
}

// AppUserSettings is app user record setup for an app
type AppUserSettings interface {
	AppUserCollectionName() string
	NewAppUserEntity() AppUser

	// AppUserEntityKind returns kind of app user entity
	// Deprecated: use AppUserEntityType() instead
	//AppUserEntityKind() string

	// AppUserEntityType returns type of a DTO struct for an app user record
	// Deprecated: use NewAppUserEntity() instead
	//AppUserEntityType() reflect.Type
}

type executionContext struct {
	c context.Context
}

func (ec executionContext) Context() context.Context {
	return ec.c
}

// NewExecutionContext creates new execution context
func NewExecutionContext(c context.Context) ExecutionContext {
	return executionContext{
		c: c,
	}
}
