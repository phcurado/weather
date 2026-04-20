package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad_Defaults_NoFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if cfg.DefaultCity != "" {
		t.Errorf("DefaultCity = %q; want empty", cfg.DefaultCity)
	}
	if cfg.Units != "metric" {
		t.Errorf("Units = %q; want metric", cfg.Units)
	}
	if cfg.CacheTTL != 10*time.Minute {
		t.Errorf("CacheTTL = %v; want 10m", cfg.CacheTTL)
	}
}

func TestLoad_ParsesFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfgDir := filepath.Join(dir, "weather")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `
default_city = "Tallinn"
units        = "imperial"
cache_ttl    = "30s"
`
	if err := os.WriteFile(filepath.Join(cfgDir, "config.toml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if cfg.DefaultCity != "Tallinn" {
		t.Errorf("DefaultCity = %q", cfg.DefaultCity)
	}
	if cfg.Units != "imperial" {
		t.Errorf("Units = %q", cfg.Units)
	}
	if cfg.CacheTTL != 30*time.Second {
		t.Errorf("CacheTTL = %v", cfg.CacheTTL)
	}
}

func TestLoad_RejectsBadUnits(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	cfgDir := filepath.Join(dir, "weather")
	_ = os.MkdirAll(cfgDir, 0o755)
	_ = os.WriteFile(filepath.Join(cfgDir, "config.toml"),
		[]byte(`units = "kelvin"`), 0o644)

	if _, err := Load(); err == nil {
		t.Fatal("expected error for bad units")
	}
}

func TestPath_UsesXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg")
	got := Path()
	want := "/tmp/xdg/weather/config.toml"
	if got != want {
		t.Errorf("Path() = %q; want %q", got, want)
	}
}
