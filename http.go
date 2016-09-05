package strongo

import (
	"net/http"
	"strings"
)

func AddHttpHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
	if !strings.HasSuffix(pattern, "/") {
		http.HandleFunc(pattern+"/", handler)
	}
}
