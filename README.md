# keep

Tiny HTTP cache proxy. Put it in front of a scraper so you stop re-fetching the same pages during dev.

```bash
go build -o keep .
KEEP_UPSTREAM=http://127.0.0.1:8000 ./keep
```

Endpoints: `/healthz`, `/readyz`, `/stats`, everything else proxies upstream with in-memory TTL cache.

Cached GET responses include `X-Keep-Cache: HIT` or `MISS` for quick debugging.

`POST /cache/purge` clears the in-memory cache (dev helper).

## Configuration

Settings load from `config.toml` when `KEEP_CONFIG` points at a file. Environment variables override file values.

```toml
upstream = "http://127.0.0.1:8000"
addr = ":8787"
ttl = "60s"
max_entries = 256
```

Environment: `KEEP_CONFIG`, `KEEP_UPSTREAM`, `KEEP_ADDR`, `KEEP_TTL` (default 60s), `KEEP_MAX_ENTRIES` (default 256). Honors upstream `Cache-Control: max-age` when present; skips cache for `no-store`.

`/stats` includes `hit_ratio` (hits / (hits + misses)).

## Docker sidecar

```bash
docker build -t keep .
docker run --rm -p 8787:8787 -e KEEP_UPSTREAM=http://host.docker.internal:8000 keep
```

## Tests

```bash
go test ./...
```

## Portfolio

Runs as a sidecar in the [`checkup`](../checkup) compose stack. Hit/miss stats appear on the `switchboard` dashboard.
