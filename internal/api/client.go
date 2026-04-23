package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ErrCityNotFound is returned when geocoding has no results.
var ErrCityNotFound = errors.New("city not found")

const (
	DefaultGeocodeBase  = "https://geocoding-api.open-meteo.com"
	DefaultForecastBase = "https://api.open-meteo.com"
)

type Options struct {
	GeocodeBase  string
	ForecastBase string
	IPGeoBase    string
	Timeout      time.Duration
}

type Client struct {
	geocodeBase  string
	forecastBase string
	ipGeoBase    string
	http         *http.Client
}

func NewClient(o Options) *Client {
	if o.GeocodeBase == "" {
		o.GeocodeBase = DefaultGeocodeBase
	}
	if o.ForecastBase == "" {
		o.ForecastBase = DefaultForecastBase
	}
	if o.IPGeoBase == "" {
		o.IPGeoBase = DefaultIPGeoBase
	}
	if o.Timeout == 0 {
		o.Timeout = 5 * time.Second
	}
	return &Client{
		geocodeBase:  o.GeocodeBase,
		forecastBase: o.ForecastBase,
		ipGeoBase:    o.IPGeoBase,
		http:         &http.Client{Timeout: o.Timeout},
	}
}

func (c *Client) getJSON(raw string, q url.Values, out any) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	u.RawQuery = q.Encode()
	resp, err := c.http.Get(u.String())
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("open-meteo %s: %s", resp.Status, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func ftoa(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }
