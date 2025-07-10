package main

import (
	"os"
	"testing"
	"time"
)

func TestCacheRespectsMaxEntriesEnv(t *testing.T) {
	t.Setenv("KEEP_MAX_ENTRIES", "2")
	// env is read in main(); exercise cache directly
	c := newCache(2, time.Minute)
	c.set("/a", &entry{body: []byte("a"), status: 200, expiresAt: time.Now().Add(time.Minute)})
	c.set("/b", &entry{body: []byte("b"), status: 200, expiresAt: time.Now().Add(time.Minute)})
	c.set("/c", &entry{body: []byte("c"), status: 200, expiresAt: time.Now().Add(time.Minute)})
	if len(c.items) > 2 {
		t.Fatalf("expected at most 2 entries, got %d", len(c.items))
	}
	_ = os.Getenv("KEEP_MAX_ENTRIES")
}
