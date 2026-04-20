package render

import (
	"fmt"

	"github.com/phcurado/weather/internal/api"
)

// Widget renders a compact single-line status for tmux. Uses Nerd Font weather
// glyphs (weather-icons set, U+E300–U+E3EB) wrapped in tmux format strings so
// tmux colorizes natively. No trailing newline.
//
// Requires a Nerd Font patched terminal font.
func Widget(w api.Weather) string {
	icon, color := widgetIcon(w.Current.WeatherCode)
	return fmt.Sprintf("#[fg=%s]%s#[default] %.0f%s",
		color, icon, w.Current.TempC, tempUnit(w.Units))
}

// widgetIcon maps a WMO code to a nerd-font glyph and a tmux color name.
func widgetIcon(code int) (string, string) {
	switch {
	case code == 0 || code == 1:
		return "\ue30d", "colour222" // wi-day-sunny, yellow
	case code == 2:
		return "\ue302", "colour222" // wi-day-cloudy, yellow
	case code == 3:
		return "\ue33d", "colour244" // wi-cloudy, gray
	case code == 45 || code == 48:
		return "\ue313", "colour244" // wi-fog, gray
	case code >= 51 && code <= 67:
		return "\ue318", "colour111" // wi-rain, blue
	case code >= 80 && code <= 82:
		return "\ue319", "colour111" // wi-showers, blue
	case (code >= 71 && code <= 77) || code == 85 || code == 86:
		return "\ue31a", "colour117" // wi-snow, cyan
	case code >= 95 && code <= 99:
		return "\ue31d", "colour215" // wi-thunderstorm, orange
	}
	return "\ue33d", "colour244"
}
