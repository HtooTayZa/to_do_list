package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type keepConfig struct {
	Upstream   string
	Addr       string
	TTL        time.Duration
	MaxEntries int
}

func loadConfig() keepConfig {
	cfg := keepConfig{
		Upstream:   "http://127.0.0.1:8000",
		Addr:       ":8787",
		TTL:        60 * time.Second,
		MaxEntries: 256,
	}

	if path := os.Getenv("KEEP_CONFIG"); path != "" {
		if data, err := os.ReadFile(path); err == nil {
			parseConfigToml(string(data), &cfg)
		}
	}

	if v := os.Getenv("KEEP_UPSTREAM"); v != "" {
		cfg.Upstream = v
	}
	if v := os.Getenv("KEEP_ADDR"); v != "" {
		cfg.Addr = v
	}
	if v := os.Getenv("KEEP_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.TTL = d
		}
	}
	if v := os.Getenv("KEEP_MAX_ENTRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MaxEntries = n
		}
	}

	return cfg
}

func parseConfigToml(raw string, cfg *keepConfig) {
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"`)
		switch key {
		case "upstream":
			cfg.Upstream = val
		case "addr":
			cfg.Addr = val
		case "ttl":
			if d, err := time.ParseDuration(val); err == nil {
				cfg.TTL = d
			}
		case "max_entries":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				cfg.MaxEntries = n
			}
		}
	}
}
