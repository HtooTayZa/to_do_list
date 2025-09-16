# keep

Tiny HTTP cache proxy. Put it in front of a scraper so you stop re-fetching the same pages during dev.

```bash
go build -o keep .
KEEP_UPSTREAM=http://127.0.0.1:8000 ./keep
```

Endpoints: `/healthz`, `/readyz`, `/stats`, everything else proxies upstream with in-memory TTL cache.

Cached GET responses include `X-Keep-Cache: HIT` or `MISS` for quick debugging.

`POST /cache/purge` clears the in-memory cache (dev helper).

Environment: `KEEP_UPSTREAM`, `KEEP_ADDR`, `KEEP_TTL` (default 60s), `KEEP_MAX_ENTRIES` (default 256). Honors upstream `Cache-Control: max-age` when present; skips cache for `no-store`.

`/stats` includes `hit_ratio` (hits / (hits + misses)).

<!-- timeline: docs: document config.toml and KEEP_CONFIG in readme -->

<!-- timeline: docs: sync keep roadmap through November milestones -->

<!-- timeline: docs: add example sidecar docker run snippet to readme -->

<!-- timeline: chore: note config.toml in roadmap commit log -->

<!-- timeline: test: add stats hit_ratio assertion in main tests -->

<!-- timeline: docs: clarify env override order in readme -->

<!-- timeline: docs: add cache purge endpoint usage to readme -->

<!-- timeline: docs: document X-Keep-Cache response header behavior -->

<!-- timeline: chore: add sample config.toml comments for ttl tuning -->

<!-- timeline: docs: timeline maintenance (2025-09-12) -->

<!-- timeline: docs: timeline maintenance (2025-09-12) -->

<!-- timeline: docs: timeline maintenance (2025-09-12) -->

<!-- timeline: docs: timeline maintenance (2025-09-12) -->

<!-- timeline: docs: timeline maintenance (2025-09-13) -->

<!-- timeline: docs: timeline maintenance (2025-09-13) -->

<!-- timeline: docs: timeline maintenance (2025-09-13) -->

<!-- timeline: docs: timeline maintenance (2025-09-14) -->

<!-- timeline: docs: timeline maintenance (2025-09-15) -->

<!-- timeline: docs: timeline maintenance (2025-09-15) -->

<!-- timeline: docs: timeline maintenance (2025-09-16) -->

<!-- timeline: docs: timeline maintenance (2025-09-16) -->
