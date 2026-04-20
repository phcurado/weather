package wmo

import "testing"

func TestLookup_KnownCodes(t *testing.T) {
	cases := []struct {
		code      int
		wantEmoji string
		wantLabel string
	}{
		{0, "☀", "Clear"},
		{1, "🌤", "Mostly clear"},
		{2, "⛅", "Partly cloudy"},
		{3, "☁", "Overcast"},
		{45, "🌫", "Fog"},
		{48, "🌫", "Rime fog"},
		{51, "🌦", "Light drizzle"},
		{61, "🌧", "Light rain"},
		{63, "🌧", "Rain"},
		{65, "🌧", "Heavy rain"},
		{71, "🌨", "Light snow"},
		{73, "🌨", "Snow"},
		{75, "🌨", "Heavy snow"},
		{80, "🌦", "Rain showers"},
		{95, "⛈", "Thunderstorm"},
		{99, "⛈", "Thunderstorm with hail"},
	}
	for _, c := range cases {
		gotEmoji, gotLabel := Lookup(c.code)
		if gotEmoji != c.wantEmoji || gotLabel != c.wantLabel {
			t.Errorf("Lookup(%d) = (%q,%q); want (%q,%q)",
				c.code, gotEmoji, gotLabel, c.wantEmoji, c.wantLabel)
		}
	}
}

func TestLookup_Unknown(t *testing.T) {
	e, l := Lookup(9999)
	if e != "?" || l != "Unknown" {
		t.Errorf("Lookup(9999) = (%q,%q); want (\"?\",\"Unknown\")", e, l)
	}
}
