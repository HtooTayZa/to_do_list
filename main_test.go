package main

import (
	"encoding/json"
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
		coalesce: newCoalescer(),
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

func TestNoStoreResponsesAreNotCached(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write([]byte("fresh"))
	}))
	defer upstream.Close()

	c := newCache(8, time.Minute)
	s := &server{
		upstream: mustParseURL(upstream.URL),
		cache:    c,
		coalesce: newCoalescer(),
		started:  time.Now(),
	}

	req := httptest.NewRequest(http.MethodGet, "/volatile", nil)
	rec := httptest.NewRecorder()
	s.handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if _, ok := c.get("/volatile"); ok {
		t.Fatal("expected no-store response not to be cached")
	}
}

func TestStatsIncludesHitRatio(t *testing.T) {
	c := newCache(4, time.Minute)
	c.hits = 3
	c.misses = 1
	s := &server{cache: c, started: time.Now()}

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rec := httptest.NewRecorder()
	s.handle(rec, req)

	var stats map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &stats); err != nil {
		t.Fatal(err)
	}
	if stats["hit_ratio"].(float64) != 0.75 {
		t.Fatalf("expected 0.75 hit_ratio, got %v", stats["hit_ratio"])
	}
}

func TestCacheHitSetsHeader(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=60")
		_, _ = w.Write([]byte("ok"))
	}))
	defer upstream.Close()

	c := newCache(8, time.Minute)
	s := &server{
		upstream: mustParseURL(upstream.URL),
		cache:    c,
		coalesce: newCoalescer(),
		started:  time.Now(),
	}

	req := httptest.NewRequest(http.MethodGet, "/page", nil)
	rec := httptest.NewRecorder()
	s.handle(rec, req)
	if rec.Header().Get("X-Keep-Cache") != "MISS" {
		t.Fatalf("first request expected MISS, got %q", rec.Header().Get("X-Keep-Cache"))
	}

	rec2 := httptest.NewRecorder()
	s.handle(rec2, req)
	if rec2.Header().Get("X-Keep-Cache") != "HIT" {
		t.Fatalf("second request expected HIT, got %q", rec2.Header().Get("X-Keep-Cache"))
	}
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}
