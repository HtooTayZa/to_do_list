package main

import "testing"

func BenchmarkCacheGet(b *testing.B) {
	c := newCache(128, 0)
	c.set("/bench", &entry{body: []byte("ok"), status: 200, expiresAt: c.ttl})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.get("/bench")
	}
}
