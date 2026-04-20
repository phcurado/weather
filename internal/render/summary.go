package render

import (
	"fmt"
	"math"
	"strings"

	"github.com/phcurado/weather/internal/api"
	"github.com/phcurado/weather/internal/wmo"
)

var compass = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

func bearing(deg int) string {
	d := float64(((deg % 360) + 360) % 360)
	idx := int(math.Round(d/45.0)) % 8
	return compass[idx]
}

func tempUnit(units string) string {
	if units == "imperial" {
		return "°F"
	}
	return "°C"
}

func windUnit(units string) string {
	if units == "imperial" {
		return "mph"
	}
	return "m/s"
}

// Summary renders the current-conditions header: title + aligned stat block.
func Summary(w api.Weather) string {
	emoji, label := wmo.Lookup(w.Current.WeatherCode)
	tu := tempUnit(w.Units)
	wu := windUnit(w.Units)

	place := w.Coords.Name
	if w.Coords.Country != "" {
		place = fmt.Sprintf("%s, %s", w.Coords.Name, w.Coords.Country)
	}

	const labelW = 11 // "Temperature"

	var b strings.Builder
	fmt.Fprintf(&b, "%s  %s  ·  %s%s%s\n\n", emoji, place, ansiDim, label, ansiReset)

	stat := func(k, v string) {
		fmt.Fprintf(&b, "%s%-*s%s  %s\n", ansiDim, labelW, k, ansiReset, v)
	}
	stat("Temperature", fmt.Sprintf("%.0f %s  (feels %.0f %s)",
		w.Current.TempC, tu, w.Current.FeelsLikeC, tu))
	stat("Wind", fmt.Sprintf("%.0f %s %s",
		w.Current.WindSpeed, wu, bearing(w.Current.WindDirection)))
	stat("Humidity", fmt.Sprintf("%d%%", w.Current.Humidity))

	return b.String()
}
