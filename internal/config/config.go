// Package config loads the weather CLI's TOML configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DefaultCity string        `toml:"default_city"`
	Units       string        `toml:"units"`
	CacheTTL    time.Duration `toml:"-"`

	RawCacheTTL string `toml:"cache_ttl"`
}

// Path returns the resolved config file path.
func Path() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "weather", "config.toml")
}

// Load reads the config file if present, else returns defaults.
func Load() (Config, error) {
	cfg := Config{
		Units:    "metric",
		CacheTTL: 10 * time.Minute,
	}

	data, err := os.ReadFile(Path())
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Units != "metric" && cfg.Units != "imperial" {
		return cfg, fmt.Errorf("invalid units %q: must be metric or imperial", cfg.Units)
	}
	if cfg.RawCacheTTL != "" {
		d, err := time.ParseDuration(cfg.RawCacheTTL)
		if err != nil {
			return cfg, fmt.Errorf("invalid cache_ttl %q: %w", cfg.RawCacheTTL, err)
		}
		cfg.CacheTTL = d
	}
	return cfg, nil
}
