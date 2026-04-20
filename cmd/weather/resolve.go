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
// warnTo receives user-facing warnings (pass io.Discard to silence).
func resolve(cityArg string, warnTo io.Writer) (resolved, error) {
	cfg, err := config.Load()
	if err != nil {
		return resolved{}, err
	}

	city := cityArg
	if city == "" {
		city = cfg.DefaultCity
	}
	if city == "" {
		return resolved{Config: cfg}, errors.New("no city provided and no default_city in config")
	}

	c := cache.New(cfg.CacheTTL)
	client := api.NewClient(api.Options{
		GeocodeBase:  os.Getenv("WEATHER_GEOCODE_BASE"),
		ForecastBase: os.Getenv("WEATHER_FORECAST_BASE"),
		Timeout:      5 * time.Second,
	})

	coords, err := c.Geocode(city)
	if errors.Is(err, cache.ErrMiss) {
		coords, err = client.Geocode(city)
		if err != nil {
			return resolved{Config: cfg}, err
		}
		_ = c.PutGeocode(city, coords)
	} else if err != nil {
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
		fmt.Fprintf(warnTo, "warning: API failed (%v); using cached data\n", apiErr)
		return resolved{Weather: wx, Config: cfg}, nil
	}
	return resolved{Config: cfg}, apiErr
}
