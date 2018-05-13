package gaestandard

import (
	"testing"

	"context"
	"github.com/strongo/app"
)

func TestGetEnvironment(t *testing.T) {

	testEnv := func(host string, expected strongo.Environment) {
		defaultVersionHostname = func(c context.Context) string {
			return host
		}
		if environment := GetEnvironment(context.Background()); environment != expected {
			t.Errorf("Unexpected environment: %v", environment)
		}
	}

	testEnv("some-app.appspot.com", strongo.EnvProduction)
	testEnv("some-app.local", strongo.EnvLocal)
}
