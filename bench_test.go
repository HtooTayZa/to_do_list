package main

import (
	"testing"
	"time"
)

func BenchmarkCacheGet(b *testing.B) {
	c := newCache(128, time.Minute)
	c.set("/bench", &entry{body: []byte("ok"), status: 200, expiresAt: time.Now().Add(time.Minute)})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.get("/bench")
	}
}
