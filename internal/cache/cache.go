// Package cache stores geocode results permanently and weather results with a TTL.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/phcurado/weather/internal/api"
)

// ErrMiss means no entry in the cache (or the file is corrupt).
var ErrMiss = errors.New("cache miss")

type Cache struct {
	TTL time.Duration
}

func New(ttl time.Duration) *Cache {
	return &Cache{TTL: ttl}
}

type envelope struct {
	FetchedAt time.Time       `json:"fetched_at"`
	Payload   json.RawMessage `json:"payload"`
}

func (c *Cache) baseDir() string {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(base, "weather")
}

func (c *Cache) pathFor(store string, keys ...any) string {
	h := sha256.New()
	for _, k := range keys {
		_, _ = fmt.Fprintf(h, "%v|", k)
	}
	return filepath.Join(c.baseDir(), store, hex.EncodeToString(h.Sum(nil))+".json")
}

func (c *Cache) write(path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	env, err := json.Marshal(envelope{FetchedAt: time.Now().UTC(), Payload: body})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, env, 0o644)
}

func (c *Cache) read(path string, out any) (time.Time, error) {
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return time.Time{}, ErrMiss
	}
	if err != nil {
		return time.Time{}, err
	}
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return time.Time{}, ErrMiss
	}
	if err := json.Unmarshal(env.Payload, out); err != nil {
		return time.Time{}, ErrMiss
	}
	return env.FetchedAt, nil
}

// Geocode returns a previously-cached coordinate for city, or ErrMiss.
func (c *Cache) Geocode(city string) (api.Coords, error) {
	var out api.Coords
	_, err := c.read(c.pathFor("geocode", city), &out)
	return out, err
}

// PutGeocode stores coords under city. Never expires.
func (c *Cache) PutGeocode(city string, coords api.Coords) error {
	return c.write(c.pathFor("geocode", city), coords)
}

// Weather returns the cached weather for coords, whether or not it is fresh.
// fresh=true means the entry is within TTL; fresh=false means stale.
func (c *Cache) Weather(coords api.Coords) (api.Weather, bool, error) {
	var out api.Weather
	fetched, err := c.read(c.pathFor("weather", coords.Lat, coords.Lon), &out)
	if err != nil {
		return api.Weather{}, false, err
	}
	fresh := time.Since(fetched) < c.TTL
	return out, fresh, nil
}

// PutWeather stores wx under coords.
func (c *Cache) PutWeather(coords api.Coords, wx api.Weather) error {
	return c.write(c.pathFor("weather", coords.Lat, coords.Lon), wx)
}
