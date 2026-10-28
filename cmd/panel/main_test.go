package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAllowedHost(t *testing.T) {
	ok := []string{
		"127.0.0.1:8080", "localhost:1", "[::1]:100", "127.0.0.1", "localhost",
		"8791-3204746d802354c7.monkeycode-ai.online",
	}
	for _, h := range ok {
		if !allowedHost(h) {
			t.Fatalf("want allow: %q", h)
		}
	}
	bad := []string{
		"", "evil.com", "evil.com:80", "169.254.169.254", "10.0.0.1:80",
		"monkeycode-ai.online", "evilmonkeycode-ai.online", "x.monkeycode-ai.online.evil.com",
	}
	for _, h := range bad {
		if allowedHost(h) {
			t.Fatalf("want reject: %q", h)
		}
	}
}

func TestLocalGuard(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := localGuard(inner)

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/api/meta", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("loopback host blocked: %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "http://evil.com/api/send", nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("foreign host allowed: %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "http://127.0.0.1/api/send", nil)
	req.Header.Set("Origin", "http://evil.com")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("foreign origin allowed: %d", rec.Code)
	}
}

func TestCommandChain(t *testing.T) {
	got, err := commandChain([]string{"event crisis.1000"}, "event crisis.1200")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, "|") != "event crisis.1000|event crisis.1200" {
		t.Fatalf("got %v", got)
	}
	if _, err := commandChain([]string{"a\nb"}, "ok"); err == nil {
		t.Fatal("newline pre should fail")
	}
	if _, err := commandChain(nil, ""); err == nil {
		t.Fatal("empty command should fail")
	}
	if _, err := commandChain([]string{"a", "b", "c", "d", "e"}, "ok"); err == nil {
		t.Fatal("too many pre commands should fail")
	}
}

func TestValidateCommand(t *testing.T) {
	ok := []string{"add_ship NAME_Battleship", " effect add_modifier = { multiplier = 999 } "}
	for _, c := range ok {
		if err := validateCommand(c); err != nil {
			t.Fatalf("want valid %q: %v", c, err)
		}
	}
	bad := []string{"", "   ", "a\nb", "a\rb"}
	for _, c := range bad {
		if err := validateCommand(c); err == nil {
			t.Fatalf("want invalid %q", c)
		}
	}
	if err := validateCommand(strings.Repeat("a", 600)); err == nil {
		t.Fatal("oversized command should be rejected")
	}
}

func TestNormalizeSettings(t *testing.T) {
	s := normalizeSettings(Settings{
		GameRoot:   ` "D:\\Stellaris" `,
		Version:    "9.9.9",
		ConsoleKey: "bogus",
	})
	if s.GameRoot != `D:\\Stellaris` {
		t.Fatalf("game root not trimmed: %q", s.GameRoot)
	}
	if s.Version != "4.4.6" {
		t.Fatalf("invalid version not reset: %q", s.Version)
	}
	if s.ConsoleKey != "oem3" {
		t.Fatalf("invalid console key not reset: %q", s.ConsoleKey)
	}
	if !validVersion("4.4.6") || validVersion("4.9.9") {
		t.Fatal("validVersion broken")
	}
}

func TestNeedsJSONBody(t *testing.T) {
	for _, p := range []string{"/api/send", "/api/open_console", "/api/settings", "/api/scan"} {
		if !needsJSONBody(p) {
			t.Fatalf("want JSON required: %s", p)
		}
	}
	if needsJSONBody("/api/meta") || needsJSONBody("/api/detect") {
		t.Fatal("bodyless endpoints must not require JSON")
	}
}

func TestSecureHeaders(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	rec := httptest.NewRecorder()
	secureHeaders(inner).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil))
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("missing nosniff")
	}
	if rec.Header().Get("Content-Security-Policy") == "" {
		t.Fatal("missing CSP")
	}
}

func TestEmbeddedWebMatchesSource(t *testing.T) {
	for _, name := range []string{"index.html", "app.css", "app.js"} {
		got, err := webFS.ReadFile("web/" + name)
		if err != nil {
			t.Fatalf("embedded %s: %v", name, err)
		}
		want, err := os.ReadFile(filepath.Join("..", "..", "web", name))
		if err != nil {
			t.Fatalf("source %s: %v", name, err)
		}
		if string(got) != string(want) {
			t.Fatalf("%s differs from web/: run cp web/* cmd/panel/web/", name)
		}
	}
}
