package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestParseMaxAge(t *testing.T) {
	d, ok := parseMaxAge("public, max-age=120")
	if !ok || d != 120*time.Second {
		t.Fatalf("got %v ok=%v", d, ok)
	}
	_, ok = parseMaxAge("no-store")
	if ok {
		t.Fatal("expected no max-age")
	}
}

func TestNon200ResponsesAreNotCached(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	}))
	defer upstream.Close()

	c := newCache(8, time.Minute)
	s := &server{
		upstream: mustParseURL(upstream.URL),
		cache:    c,
		started:  time.Now(),
	}

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	s.handle(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status %d", rec.Code)
	}
	if _, ok := c.get("/missing"); ok {
		t.Fatal("expected 404 response not to be cached")
	}
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}
