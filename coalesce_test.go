package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestCoalescedMissFetchesOnce(t *testing.T) {
	var hits int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Cache-Control", "max-age=60")
		_, _ = w.Write([]byte("payload"))
	}))
	defer upstream.Close()

	c := newCache(8, time.Minute)
	s := &server{
		upstream: mustParseURL(upstream.URL),
		cache:    c,
		coalesce: newCoalescer(),
		started:  time.Now(),
	}

	const n = 8
	done := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/page", nil)
			rec := httptest.NewRecorder()
			s.handle(rec, req)
			if rec.Code != http.StatusOK || rec.Body.String() != "payload" {
				t.Errorf("bad response: %d %q", rec.Code, rec.Body.String())
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < n; i++ {
		<-done
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("expected 1 upstream hit, got %d", got)
	}
}
