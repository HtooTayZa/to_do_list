package main

import "sync"

// coalescer deduplicates in-flight cache miss fetches for the same key.
type coalescer struct {
	mu      sync.Mutex
	flights map[string]*flight
}

type flight struct {
	done  chan struct{}
	entry *entry
	err   error
}

func newCoalescer() *coalescer {
	return &coalescer{flights: make(map[string]*flight)}
}

func (c *coalescer) do(key string, fn func() (*entry, error)) (*entry, error) {
	c.mu.Lock()
	if f, ok := c.flights[key]; ok {
		c.mu.Unlock()
		<-f.done
		return f.entry, f.err
	}
	f := &flight{done: make(chan struct{})}
	c.flights[key] = f
	c.mu.Unlock()

	f.entry, f.err = fn()

	c.mu.Lock()
	delete(c.flights, key)
	close(f.done)
	c.mu.Unlock()
	return f.entry, f.err
}
