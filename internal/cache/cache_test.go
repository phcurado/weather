package cache

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/phcurado/weather/internal/api"
)

func newCache(t *testing.T) *Cache {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", dir)
	return New(10 * time.Minute)
}

func TestGeocode_MissThenHit(t *testing.T) {
	c := newCache(t)

	if _, err := c.Geocode("Tallinn"); !errors.Is(err, ErrMiss) {
		t.Fatalf("first lookup err = %v; want ErrMiss", err)
	}
	want := api.Coords{Name: "Tallinn", Lat: 59.44, Lon: 24.75}
	if err := c.PutGeocode("Tallinn", want); err != nil {
		t.Fatal(err)
	}
	got, err := c.Geocode("Tallinn")
	if err != nil {
		t.Fatalf("Geocode after put err = %v", err)
	}
	if got != want {
		t.Errorf("got = %+v; want %+v", got, want)
	}
}

func TestWeather_TTLExpiry(t *testing.T) {
	c := newCache(t)
	c.TTL = 50 * time.Millisecond

	wx := api.Weather{Units: "metric", FetchedAt: time.Now().UTC()}
	if err := c.PutWeather(api.Coords{Lat: 1, Lon: 2}, wx); err != nil {
		t.Fatal(err)
	}
	got, fresh, err := c.Weather(api.Coords{Lat: 1, Lon: 2})
	if err != nil || !fresh {
		t.Fatalf("fresh read: err=%v fresh=%v", err, fresh)
	}
	if got.Units != "metric" {
		t.Errorf("got Units = %q", got.Units)
	}

	time.Sleep(80 * time.Millisecond)

	_, fresh, err = c.Weather(api.Coords{Lat: 1, Lon: 2})
	if err != nil {
		t.Fatalf("stale read err = %v", err)
	}
	if fresh {
		t.Error("expected stale")
	}
}

func TestWeather_Miss(t *testing.T) {
	c := newCache(t)
	if _, _, err := c.Weather(api.Coords{Lat: 9, Lon: 9}); !errors.Is(err, ErrMiss) {
		t.Fatalf("err = %v; want ErrMiss", err)
	}
}

func TestPutWeather_PersistsJSON(t *testing.T) {
	c := newCache(t)
	coords := api.Coords{Lat: 1, Lon: 2}
	wx := api.Weather{Units: "metric", FetchedAt: time.Now()}
	if err := c.PutWeather(coords, wx); err != nil {
		t.Fatal(err)
	}

	path := c.pathFor("weather", coords.Lat, coords.Lon)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatal(err)
	}
	if _, ok := envelope["fetched_at"]; !ok {
		t.Error("missing fetched_at")
	}
	if _, ok := envelope["payload"]; !ok {
		t.Error("missing payload")
	}
}
