package strongoapp

import (
	"context"
	"net/http"
	"strings"
)

var _ HttpAppHost = (*DefaultHttpAppHost)(nil)

var LocalHostEnv = "local"
var UnknownEnv = "unknown"

type DefaultHttpAppHost struct {
}

func (d DefaultHttpAppHost) GetEnvironment(_ context.Context, r *http.Request) string {
	if r.Host == "localhost" || strings.HasPrefix(r.Host, "localhost:") {
		return LocalHostEnv
	}
	return UnknownEnv
}

// HandleWithContext calls handler with a context.Background()
func (d DefaultHttpAppHost) HandleWithContext(handler HttpHandlerWithContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(context.Background(), w, r)
	}
}
