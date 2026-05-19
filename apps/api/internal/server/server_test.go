package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"euromillones/internal/platform/config"
)

func TestCORSAllowsConfiguredOrigin(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://euromillones-pied.vercel.app")
	t.Setenv("DRAWS_DATA_DIR", t.TempDir())
	config.Load()

	app, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodOptions, "/api/stats/hot-cold", nil)
	req.Header.Set("Origin", "https://euromillones-pied.vercel.app")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if got, want := resp.Header.Get("Access-Control-Allow-Origin"), "https://euromillones-pied.vercel.app"; got != want {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, want)
	}
}

func TestCORSAllowsFrontendURLWithoutScheme(t *testing.T) {
	t.Setenv("FRONTEND_URL", "euromillones-pied.vercel.app/")
	t.Setenv("DRAWS_DATA_DIR", t.TempDir())
	config.Load()

	app, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodOptions, "/api/stats/frequencies", nil)
	req.Header.Set("Origin", "https://euromillones-pied.vercel.app")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	req.Header.Set("Access-Control-Request-Headers", "content-type")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	if got, want := resp.Header.Get("Access-Control-Allow-Origin"), "https://euromillones-pied.vercel.app"; got != want {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, want)
	}
}
