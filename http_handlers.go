package strongo

import (
	"net/http"
	"context"
)

type ContextHandler func(c context.Context, w http.ResponseWriter, r *http.Request) // TODO: Should be somewhere else?
type HandleWithContext func(handler ContextHandler) func(w http.ResponseWriter, r *http.Request)
