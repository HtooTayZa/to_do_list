# keep

Tiny HTTP cache proxy. Put it in front of a scraper so you stop re-fetching the same pages during dev.

```bash
go build -o keep .
KEEP_UPSTREAM=http://127.0.0.1:8000 ./keep
```

Endpoints: `/healthz`, `/stats`, everything else proxies upstream with in-memory TTL cache.

Config: `config.toml` or env `KEEP_UPSTREAM`, `KEEP_ADDR`, `KEEP_TTL`, `KEEP_MAX_ENTRIES`.
