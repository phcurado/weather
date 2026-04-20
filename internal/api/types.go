// Package api speaks the Open-Meteo forecast and geocoding APIs.
package api

import "time"

// Coords is a geocoded location.
type Coords struct {
	Name     string  `json:"name"`
	Country  string  `json:"country"`
	Lat      float64 `json:"latitude"`
	Lon      float64 `json:"longitude"`
	Timezone string  `json:"timezone"`
}

// Current is instantaneous weather.
type Current struct {
	TempC         float64 `json:"temperature_2m"`
	FeelsLikeC    float64 `json:"apparent_temperature"`
	Humidity      int     `json:"relative_humidity_2m"`
	WindSpeed     float64 `json:"wind_speed_10m"`
	WindDirection int     `json:"wind_direction_10m"`
	WeatherCode   int     `json:"weather_code"`
}

// Daily is a single day of forecast.
type Daily struct {
	Date        time.Time `json:"date"`
	WeatherCode int       `json:"weather_code"`
	TempMaxC    float64   `json:"temperature_2m_max"`
	TempMinC    float64   `json:"temperature_2m_min"`
	PrecipMM    float64   `json:"precipitation_sum"`
}

// Hourly is a single hour of forecast.
type Hourly struct {
	Time        time.Time `json:"time"`
	WeatherCode int       `json:"weather_code"`
	TempC       float64   `json:"temperature_2m"`
	PrecipMM    float64   `json:"precipitation"`
}

// Weather is everything a render package needs.
type Weather struct {
	Coords    Coords
	Units     string
	Current   Current
	Daily     []Daily
	Hourly    []Hourly
	FetchedAt time.Time
}
