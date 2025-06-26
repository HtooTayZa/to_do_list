package main

import (
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
