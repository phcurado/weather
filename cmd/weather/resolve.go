package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/phcurado/weather/internal/api"
	"github.com/phcurado/weather/internal/cache"
	"github.com/phcurado/weather/internal/config"
)

type resolved struct {
	Weather api.Weather
	Config  config.Config
}

// resolve loads config, picks a city, and returns weather — using cache when
// possible and falling back to a stale cache entry on API failure.
//
// Resolution order for location:
//  1. cityArg (explicit)
//  2. IP geolocation (default — tracks travel)
//  3. cfg.City (fallback when IP lookup fails)
//
// IP geolocation results are not cached, so they follow you as you move.
// warnTo receives user-facing warnings (pass io.Discard to silence).
func resolve(cityArg string, warnTo io.Writer) (resolved, error) {
	cfg, err := config.Load()
	if err != nil {
		return resolved{}, err
	}

	c := cache.New(cfg.CacheTTL)
	client := api.NewClient(api.Options{
		GeocodeBase:  os.Getenv("WEATHER_GEOCODE_BASE"),
		ForecastBase: os.Getenv("WEATHER_FORECAST_BASE"),
		IPGeoBase:    os.Getenv("WEATHER_IPGEO_BASE"),
		Timeout:      5 * time.Second,
	})

	coords, err := resolveCoords(cityArg, cfg, client, c, warnTo)
	if err != nil {
		return resolved{Config: cfg}, err
	}

	wx, fresh, cacheErr := c.Weather(coords)
	if cacheErr == nil && fresh {
		return resolved{Weather: wx, Config: cfg}, nil
	}

	fetched, apiErr := client.Forecast(coords, cfg.Units)
	if apiErr == nil {
		_ = c.PutWeather(coords, fetched)
		return resolved{Weather: fetched, Config: cfg}, nil
	}

	if cacheErr == nil {
		_, _ = fmt.Fprintf(warnTo, "warning: API failed (%v); using cached data\n", apiErr)
		return resolved{Weather: wx, Config: cfg}, nil
	}
	return resolved{Config: cfg}, apiErr
}

func resolveCoords(cityArg string, cfg config.Config, client *api.Client, c *cache.Cache, warnTo io.Writer) (api.Coords, error) {
	if cityArg != "" {
		return geocodeCity(cityArg, client, c)
	}
	coords, ipErr := client.LocateByIP()
	if ipErr == nil {
		return coords, nil
	}
	if cfg.City != "" {
		_, _ = fmt.Fprintf(warnTo, "warning: IP geolocation failed (%v); falling back to configured city %q\n", ipErr, cfg.City)
		return geocodeCity(cfg.City, client, c)
	}
	return api.Coords{}, fmt.Errorf("locate by IP: %w (set `city` in config as a fallback or pass a city arg)", ipErr)
}

func geocodeCity(city string, client *api.Client, c *cache.Cache) (api.Coords, error) {
	coords, err := c.Geocode(city)
	if errors.Is(err, cache.ErrMiss) {
		coords, err = client.Geocode(city)
		if err != nil {
			return api.Coords{}, err
		}
		_ = c.PutGeocode(city, coords)
		return coords, nil
	}
	if err != nil {
		return api.Coords{}, err
	}
	return coords, nil
}
