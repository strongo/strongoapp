package strongo

import (
	"net/http"
	"strings"
)

// AddHTTPHandler adds http handler with / suffix
// TODO: mark as deprecated?
func AddHTTPHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
	if !strings.HasSuffix(pattern, "/") {
		http.HandleFunc(pattern+"/", handler)
	}
}
