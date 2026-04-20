package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const geocodeJSON = `{
	"results":[{
		"name":"Tallinn","country":"Estonia",
		"latitude":59.44,"longitude":24.75,"timezone":"Europe/Tallinn"
	}]
}`

const geocodeEmptyJSON = `{"generationtime_ms":0.12}`

const forecastJSON = `{
	"current":{
		"temperature_2m":12.3,"apparent_temperature":10.1,
		"relative_humidity_2m":62,"wind_speed_10m":4.0,
		"wind_direction_10m":315,"weather_code":0
	},
	"daily":{
		"time":["2026-04-18","2026-04-19","2026-04-20","2026-04-21","2026-04-22","2026-04-23","2026-04-24"],
		"weather_code":[0,61,2,0,0,61,2],
		"temperature_2m_max":[18,14,16,19,20,15,17],
		"temperature_2m_min":[8,9,7,8,10,9,8],
		"precipitation_sum":[0,5,1,0,0,8,2]
	}
}`

func TestGeocode_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v1/search") {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if r.URL.Query().Get("name") != "Tallinn" {
			t.Errorf("missing name query")
		}
		_, _ = w.Write([]byte(geocodeJSON))
	}))
	defer srv.Close()

	c := NewClient(Options{GeocodeBase: srv.URL, ForecastBase: "unused", Timeout: 2 * time.Second})
	got, err := c.Geocode("Tallinn")
	if err != nil {
		t.Fatalf("Geocode err = %v", err)
	}
	if got.Name != "Tallinn" || got.Lat != 59.44 || got.Lon != 24.75 {
		t.Errorf("got = %+v", got)
	}
}

func TestGeocode_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(geocodeEmptyJSON))
	}))
	defer srv.Close()

	c := NewClient(Options{GeocodeBase: srv.URL, ForecastBase: "unused", Timeout: 2 * time.Second})
	_, err := c.Geocode("Nowhereville")
	if err == nil {
		t.Fatal("expected ErrCityNotFound")
	}
}

func TestForecast_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("latitude") != "59.44" || q.Get("longitude") != "24.75" {
			t.Errorf("bad coords in query: %v", q)
		}
		if q.Get("temperature_unit") != "celsius" {
			t.Errorf("units = %q; want celsius", q.Get("temperature_unit"))
		}
		_, _ = w.Write([]byte(forecastJSON))
	}))
	defer srv.Close()

	c := NewClient(Options{GeocodeBase: "unused", ForecastBase: srv.URL, Timeout: 2 * time.Second})
	coords := Coords{Name: "Tallinn", Lat: 59.44, Lon: 24.75, Timezone: "Europe/Tallinn"}
	wx, err := c.Forecast(coords, "metric")
	if err != nil {
		t.Fatalf("Forecast err = %v", err)
	}
	if wx.Current.TempC != 12.3 {
		t.Errorf("TempC = %v", wx.Current.TempC)
	}
	if len(wx.Daily) != 7 {
		t.Fatalf("len(Daily) = %d; want 7", len(wx.Daily))
	}
	if wx.Daily[1].PrecipMM != 5 || wx.Daily[1].WeatherCode != 61 {
		t.Errorf("Daily[1] = %+v", wx.Daily[1])
	}
}

func TestForecast_Imperial_UsesFahrenheit(t *testing.T) {
	gotUnit := ""
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUnit = r.URL.Query().Get("temperature_unit")
		_, _ = w.Write([]byte(forecastJSON))
	}))
	defer srv.Close()

	c := NewClient(Options{GeocodeBase: "unused", ForecastBase: srv.URL, Timeout: 2 * time.Second})
	_, err := c.Forecast(Coords{Lat: 1, Lon: 1}, "imperial")
	if err != nil {
		t.Fatalf("Forecast err = %v", err)
	}
	if gotUnit != "fahrenheit" {
		t.Errorf("temperature_unit = %q; want fahrenheit", gotUnit)
	}
}

func TestForecast_5xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", 500)
	}))
	defer srv.Close()

	c := NewClient(Options{GeocodeBase: "unused", ForecastBase: srv.URL, Timeout: 2 * time.Second})
	_, err := c.Forecast(Coords{Lat: 1, Lon: 1}, "metric")
	if err == nil {
		t.Fatal("expected error on 500")
	}
}
