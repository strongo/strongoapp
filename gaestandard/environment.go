package gaestandard

import (
	"context"
	"github.com/strongo/app"
	"google.golang.org/appengine"
	"strings"
)

var defaultVersionHostname = appengine.DefaultVersionHostname

func GetEnvironment(c context.Context) strongo.Environment {
	hostname := defaultVersionHostname(c)
	return GetEnvironmentFromHost(hostname)
}

func GetEnvironmentFromHost(host string) strongo.Environment {
	if strings.Contains(host, "dev") && strings.HasSuffix(host, ".appspot.com") {
		return strongo.EnvDevTest
	} else if host == "localhost" || strings.HasPrefix(host, "localhost:") || strings.HasSuffix(host, ".ngrok.io") || strings.Contains(host, "local") {
		return strongo.EnvLocal
	} else {
		return strongo.EnvProduction
	}
	return strongo.EnvUnknown
}
