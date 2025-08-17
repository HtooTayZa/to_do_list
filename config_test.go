package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfigFromToml(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `# test
upstream = "http://cache-upstream:9000"
addr = ":9999"
ttl = "2m"
max_entries = 64
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("KEEP_CONFIG", path)
	t.Setenv("KEEP_UPSTREAM", "")
	t.Setenv("KEEP_ADDR", "")
	t.Setenv("KEEP_TTL", "")
	t.Setenv("KEEP_MAX_ENTRIES", "")

	cfg := loadConfig()
	if cfg.Upstream != "http://cache-upstream:9000" {
		t.Fatalf("upstream %q", cfg.Upstream)
	}
	if cfg.Addr != ":9999" {
		t.Fatalf("addr %q", cfg.Addr)
	}
	if cfg.TTL != 2*time.Minute {
		t.Fatalf("ttl %v", cfg.TTL)
	}
	if cfg.MaxEntries != 64 {
		t.Fatalf("max_entries %d", cfg.MaxEntries)
	}
}

func TestEnvOverridesToml(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(`upstream = "http://from-file"`), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("KEEP_CONFIG", path)
	t.Setenv("KEEP_UPSTREAM", "http://from-env")

	cfg := loadConfig()
	if cfg.Upstream != "http://from-env" {
		t.Fatalf("expected env override, got %q", cfg.Upstream)
	}
}
