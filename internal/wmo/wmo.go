// Package wmo maps WMO weather interpretation codes to an emoji + label.
// Codes per Open-Meteo docs: https://open-meteo.com/en/docs
package wmo

type entry struct {
	emoji string
	label string
}

var table = map[int]entry{
	0:  {"☀", "Clear"},
	1:  {"🌤", "Mostly clear"},
	2:  {"⛅", "Partly cloudy"},
	3:  {"☁", "Overcast"},
	45: {"🌫", "Fog"},
	48: {"🌫", "Rime fog"},
	51: {"🌦", "Light drizzle"},
	53: {"🌦", "Drizzle"},
	55: {"🌦", "Heavy drizzle"},
	56: {"🌧", "Freezing drizzle"},
	57: {"🌧", "Heavy freezing drizzle"},
	61: {"🌧", "Light rain"},
	63: {"🌧", "Rain"},
	65: {"🌧", "Heavy rain"},
	66: {"🌧", "Freezing rain"},
	67: {"🌧", "Heavy freezing rain"},
	71: {"🌨", "Light snow"},
	73: {"🌨", "Snow"},
	75: {"🌨", "Heavy snow"},
	77: {"🌨", "Snow grains"},
	80: {"🌦", "Rain showers"},
	81: {"🌧", "Heavy rain showers"},
	82: {"🌧", "Violent rain showers"},
	85: {"🌨", "Snow showers"},
	86: {"🌨", "Heavy snow showers"},
	95: {"⛈", "Thunderstorm"},
	96: {"⛈", "Thunderstorm with hail"},
	99: {"⛈", "Thunderstorm with hail"},
}

// Lookup returns emoji and human label for a WMO weather code.
// Unknown codes return ("?", "Unknown").
func Lookup(code int) (string, string) {
	if e, ok := table[code]; ok {
		return e.emoji, e.label
	}
	return "?", "Unknown"
}
