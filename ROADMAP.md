# keep — ROADMAP

Tiny HTTP caching reverse proxy with TTL and `/stats`.

**Active window:** Sep 12 – Nov 28, 2025

---

## Milestone 1 — Reverse proxy skeleton

- [x] `net/http` server with configurable upstream (`KEEP_UPSTREAM`)
- [x] Non-GET requests pass through reverse proxy

## Milestone 2 — TTL cache with size cap

- [x] In-memory cache with TTL eviction
- [x] `KEEP_MAX_ENTRIES` size cap with eviction counter
- [x] Honor `Cache-Control: max-age` on cacheable GET responses

## Milestone 3 — Stampede protection on miss

- [x] Request coalescing (singleflight) on concurrent cache misses

## Milestone 4 — Stats + health endpoints

- [x] `/healthz`, `/readyz`, `/stats` JSON (hits, misses, evictions, uptime)

## Milestone 5 — Config + graceful shutdown

- [x] Env-based config (`KEEP_ADDR`, `KEEP_TTL`, etc.)
- [x] SIGINT/SIGTERM graceful shutdown

## Milestone 6 — Dockerfile + README

- [x] Dockerfile for sidecar deployment
- [x] Short README with env reference

---

## Commit log

| Simulated date | Commit | Milestone |
|----------------|--------|-----------|
| 2025-09-12 | chore: initialize keep repo with go module layout | setup |
| 2025-09-15 | feat: add reverse proxy skeleton with health endpoints | 1, 4 |
| 2025-09-20 | feat: add in-memory TTL cache with size cap | 2 |
| 2025-10-02 | feat: honor Cache-Control max-age when caching GET responses | 2 |
| 2025-10-10 | test: assert 404 responses are not cached | 2 |
| 2025-10-18 | feat: support KEEP_MAX_ENTRIES env for cache size | 2 |
| 2025-11-05 | docs: align readme with env-based config | 5 |
| 2025-11-12 | chore: add Dockerfile for keep sidecar | 6 |
| 2025-11-20 | feat: coalesce concurrent cache miss fetches | 3 |
| 2025-11-25 | fix: use graceful shutdown with timeout on signals | 5 |
| 2025-08-15 | feat: add X-Keep-Cache hit/miss response header | 4 |
