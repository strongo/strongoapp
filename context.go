package strongoapp

import (
	"context"
	"github.com/strongo/app/appuser"
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

	// AppUserCollectionName returns name of a collection for app user records
	AppUserCollectionName() string // TODO: Add a link to example of usage

	// NewAppUserData TODO: Needs documentation on intended use and examples of usage
	NewAppUserData() appuser.BaseUserData // TODO: Consider returning dalgo record
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
