package gaestandard

import (
	"github.com/strongo/app"
	"google.golang.org/appengine/v2"
	"net/http"
)

func HandleWithContext(handler strongo.ContextHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		handler(c, w, r)
	}
}
