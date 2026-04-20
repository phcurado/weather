package api

import (
	"fmt"
	"net/url"
	"time"
)

type forecastResp struct {
	Current Current `json:"current"`
	Daily   struct {
		Time        []string  `json:"time"`
		WeatherCode []int     `json:"weather_code"`
		TempMax     []float64 `json:"temperature_2m_max"`
		TempMin     []float64 `json:"temperature_2m_min"`
		Precip      []float64 `json:"precipitation_sum"`
	} `json:"daily"`
	Hourly struct {
		Time        []string  `json:"time"`
		WeatherCode []int     `json:"weather_code"`
		Temp        []float64 `json:"temperature_2m"`
		Precip      []float64 `json:"precipitation"`
	} `json:"hourly"`
}

// Forecast fetches current + 7-day daily forecast for coords.
// units is "metric" or "imperial".
func (c *Client) Forecast(coords Coords, units string) (Weather, error) {
	q := url.Values{}
	q.Set("latitude", ftoa(coords.Lat))
	q.Set("longitude", ftoa(coords.Lon))
	q.Set("current", "temperature_2m,apparent_temperature,relative_humidity_2m,wind_speed_10m,wind_direction_10m,weather_code")
	q.Set("daily", "weather_code,temperature_2m_max,temperature_2m_min,precipitation_sum")
	q.Set("hourly", "weather_code,temperature_2m,precipitation")
	q.Set("forecast_days", "7")
	q.Set("timezone", "auto")
	switch units {
	case "imperial":
		q.Set("temperature_unit", "fahrenheit")
		q.Set("wind_speed_unit", "mph")
		q.Set("precipitation_unit", "inch")
	default:
		q.Set("temperature_unit", "celsius")
		q.Set("wind_speed_unit", "ms")
		q.Set("precipitation_unit", "mm")
	}

	var out forecastResp
	if err := c.getJSON(c.forecastBase+"/v1/forecast", q, &out); err != nil {
		return Weather{}, err
	}

	wx := Weather{
		Coords:    coords,
		Units:     units,
		Current:   out.Current,
		FetchedAt: time.Now().UTC(),
	}
	for i := range out.Daily.Time {
		t, err := time.Parse("2006-01-02", out.Daily.Time[i])
		if err != nil {
			return Weather{}, fmt.Errorf("parse daily time %q: %w", out.Daily.Time[i], err)
		}
		wx.Daily = append(wx.Daily, Daily{
			Date:        t,
			WeatherCode: out.Daily.WeatherCode[i],
			TempMaxC:    out.Daily.TempMax[i],
			TempMinC:    out.Daily.TempMin[i],
			PrecipMM:    out.Daily.Precip[i],
		})
	}
	for i := range out.Hourly.Time {
		t, err := time.Parse("2006-01-02T15:04", out.Hourly.Time[i])
		if err != nil {
			return Weather{}, fmt.Errorf("parse hourly time %q: %w", out.Hourly.Time[i], err)
		}
		wx.Hourly = append(wx.Hourly, Hourly{
			Time:        t,
			WeatherCode: out.Hourly.WeatherCode[i],
			TempC:       out.Hourly.Temp[i],
			PrecipMM:    out.Hourly.Precip[i],
		})
	}
	return wx, nil
}
