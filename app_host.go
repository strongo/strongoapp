package strongo

import (
	"context"
	"net/http"
)

type HttpHandlerWithContext func(c context.Context, w http.ResponseWriter, r *http.Request) // TODO: Should be somewhere else?
//type HandleWithContext func(handler HttpHandlerWithContext) func(w http.ResponseWriter, r *http.Request)

type HttpAppHost interface {
	GetEnvironment(c context.Context, r *http.Request) Environment
	HandleWithContext(handler HttpHandlerWithContext) func(w http.ResponseWriter, r *http.Request)
}
