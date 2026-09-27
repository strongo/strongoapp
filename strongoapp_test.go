package strongoapp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testContextKey struct{}

func TestExecutionContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), testContextKey{}, "val")
	ec := NewExecutionContext(ctx)
	if ec.Context() != ctx {
		t.Errorf("expected %v, got %v", ctx, ec.Context())
	}
}

func TestDefaultHttpAppHost_GetEnvironment(t *testing.T) {
	host := DefaultHttpAppHost{}

	reqLocal := httptest.NewRequest(http.MethodGet, "http://localhost/", nil)
	if env := host.GetEnvironment(context.Background(), reqLocal); env != LocalHostEnv {
		t.Errorf("expected %s, got %s", LocalHostEnv, env)
	}

	reqLocalPort := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	if env := host.GetEnvironment(context.Background(), reqLocalPort); env != LocalHostEnv {
		t.Errorf("expected %s, got %s", LocalHostEnv, env)
	}

	reqRemote := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
	if env := host.GetEnvironment(context.Background(), reqRemote); env != UnknownEnv {
		t.Errorf("expected %s, got %s", UnknownEnv, env)
	}
}

func TestDefaultHttpAppHost_HandleWithContext(t *testing.T) {
	host := DefaultHttpAppHost{}
	called := false
	handler := host.HandleWithContext(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		called = true
		if c == nil {
			t.Error("context should not be nil")
		}
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	handler(rec, req)

	if !called {
		t.Error("handler was not called")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAddHTTPHandler(t *testing.T) {
	called := 0
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		called++
	}

	// Pattern without trailing slash (registers both /test-no-slash and /test-no-slash/)
	AddHTTPHandler("/test-no-slash", dummyHandler)

	// Pattern with trailing slash (registers only /test-with-slash/)
	AddHTTPHandler("/test-with-slash/", dummyHandler)
}
