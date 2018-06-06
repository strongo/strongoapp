package gaestandard

import (
	"net/http"
	"google.golang.org/appengine"
	"github.com/strongo/app"
)

func HandleWithContext(handler strongo.ContextHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		handler(c, w, r)
	}
}
