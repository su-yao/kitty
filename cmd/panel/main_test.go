package main

import (
	"net/http"
	"net/http/httptest"
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
