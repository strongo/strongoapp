package strongoapp

import (
	"context"
	"net/http"
)

// HttpHandlerWithContext - TODO: document purpose
type HttpHandlerWithContext func(c context.Context, w http.ResponseWriter, r *http.Request)

// HttpAppHost - TODO: document purpose
type HttpAppHost interface {

	// GetEnvironment determines environment based on request
	GetEnvironment(c context.Context, r *http.Request) string

	// HandleWithContext - calls handler with a context.Context specific to app/host/request
	HandleWithContext(handler HttpHandlerWithContext) func(w http.ResponseWriter, r *http.Request)
}
