package strongoapp

import (
	"context"
	"net/http"
	"strings"
)

type HttpHandlerWithContext func(c context.Context, w http.ResponseWriter, r *http.Request) // TODO: Should be somewhere else?
//type HandleWithContext func(handler HttpHandlerWithContext) func(w http.ResponseWriter, r *http.Request)

type HttpAppHost interface {
	GetEnvironment(c context.Context, r *http.Request) Environment
	HandleWithContext(handler HttpHandlerWithContext) func(w http.ResponseWriter, r *http.Request)
}

var _ HttpAppHost = (*DefaultHttpAppHost)(nil)

type DefaultHttpAppHost struct {
}

func (d DefaultHttpAppHost) GetEnvironment(c context.Context, r *http.Request) Environment {
	if r.Host == "localhost" || strings.HasPrefix(r.Host, "localhost:") {
		return EnvLocal
	}
	return EnvUnknown
}

func (d DefaultHttpAppHost) HandleWithContext(handler HttpHandlerWithContext) func(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
