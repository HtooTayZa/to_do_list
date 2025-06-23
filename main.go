package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type entry struct {
	body      []byte
	status    int
	header    http.Header
	expiresAt time.Time
}

type cache struct {
	mu        sync.RWMutex
	items     map[string]*entry
	max       int
	ttl       time.Duration
	hits      int64
	misses    int64
	evictions int64
}

func newCache(max int, ttl time.Duration) *cache {
	return &cache{items: make(map[string]*entry), max: max, ttl: ttl}
}

func (c *cache) get(key string) (*entry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[key]
	if !ok || time.Now().After(e.expiresAt) {
		if ok {
			delete(c.items, key)
		}
		c.misses++
		return nil, false
	}
	c.hits++
	return e, true
}

func (c *cache) set(key string, e *entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.items) >= c.max {
		for k := range c.items {
			delete(c.items, k)
			c.evictions++
			break
		}
	}
	c.items[key] = e
}

type server struct {
	upstream *url.URL
	proxy    *httputil.ReverseProxy
	cache    *cache
	started  time.Time
}

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
		return
	}
	if r.URL.Path == "/stats" {
		s.cache.mu.RLock()
		stats := map[string]any{
			"hits":      s.cache.hits,
			"misses":    s.cache.misses,
			"evictions": s.cache.evictions,
			"entries":   len(s.cache.items),
			"uptime_s":  int(time.Since(s.started).Seconds()),
		}
		s.cache.mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(stats)
		return
	}

	if r.Method != http.MethodGet {
		s.proxy.ServeHTTP(w, r)
		return
	}

	key := r.URL.String()
	if e, ok := s.cache.get(key); ok {
		copyHeader(w.Header(), e.header)
		w.WriteHeader(e.status)
		_, _ = w.Write(e.body)
		return
	}

	target := s.upstream.ResolveReference(r.URL).String()
	resp, err := http.Get(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	e := &entry{
		body:      body,
		status:    resp.StatusCode,
		header:    resp.Header.Clone(),
		expiresAt: time.Now().Add(s.cache.ttl),
	}
	s.cache.set(key, e)
	copyHeader(w.Header(), e.header)
	w.WriteHeader(e.status)
	_, _ = w.Write(e.body)
}

func copyHeader(dst, src http.Header) {
	for k, vals := range src {
		for _, v := range vals {
			dst.Add(k, v)
		}
	}
}

func main() {
	upstream := os.Getenv("KEEP_UPSTREAM")
	if upstream == "" {
		upstream = "http://127.0.0.1:8000"
	}
	addr := os.Getenv("KEEP_ADDR")
	if addr == "" {
		addr = ":8787"
	}
	ttl := 60 * time.Second
	if v := os.Getenv("KEEP_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			ttl = d
		}
	}

	u, err := url.Parse(upstream)
	if err != nil {
		slog.Error("bad upstream", "err", err)
		os.Exit(1)
	}

	s := &server{
		upstream: u,
		proxy:    httputil.NewSingleHostReverseProxy(u),
		cache:    newCache(256, ttl),
		started:  time.Now(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)

	httpServer := &http.Server{Addr: addr, Handler: mux}
	go func() {
		slog.Info("keep listening", "addr", addr, "upstream", upstream)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	_ = httpServer.Close()
}
