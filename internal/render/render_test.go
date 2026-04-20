package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/phcurado/weather/internal/api"
)

func init() {
	DisableColor()
}

func fixtureWeather(units string) api.Weather {
	return api.Weather{
		Coords:    api.Coords{Name: "Tallinn", Country: "Estonia", Lat: 59.44, Lon: 24.75, Timezone: "Europe/Tallinn"},
		Units:     units,
		FetchedAt: time.Date(2026, 4, 18, 10, 0, 0, 0, time.UTC),
		Current: api.Current{
			TempC: 12, FeelsLikeC: 10, Humidity: 62,
			WindSpeed: 4, WindDirection: 315, WeatherCode: 0,
		},
		Daily: []api.Daily{
			{Date: time.Date(2026, 4, 18, 0, 0, 0, 0, time.UTC), WeatherCode: 0, TempMaxC: 18, TempMinC: 8, PrecipMM: 0},
			{Date: time.Date(2026, 4, 19, 0, 0, 0, 0, time.UTC), WeatherCode: 61, TempMaxC: 14, TempMinC: 9, PrecipMM: 5},
			{Date: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC), WeatherCode: 2, TempMaxC: 16, TempMinC: 7, PrecipMM: 1},
			{Date: time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC), WeatherCode: 0, TempMaxC: 19, TempMinC: 8, PrecipMM: 0},
			{Date: time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC), WeatherCode: 0, TempMaxC: 20, TempMinC: 10, PrecipMM: 0},
			{Date: time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC), WeatherCode: 61, TempMaxC: 15, TempMinC: 9, PrecipMM: 8},
			{Date: time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC), WeatherCode: 2, TempMaxC: 17, TempMinC: 8, PrecipMM: 2},
		},
	}
}

func goldenAssert(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		_ = os.WriteFile(path, []byte(got), 0o644)
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if strings.TrimRight(got, "\n") != strings.TrimRight(string(want), "\n") {
		t.Errorf("mismatch for %s\n--- want ---\n%s\n--- got ---\n%s", name, want, got)
	}
}

func TestSummary(t *testing.T) {
	got := Summary(fixtureWeather("metric"))
	goldenAssert(t, "summary.txt", got)
}

func TestForecast(t *testing.T) {
	got := Forecast(fixtureWeather("metric"))
	goldenAssert(t, "forecast.txt", got)
}

func TestWidget_Metric(t *testing.T) {
	got := Widget(fixtureWeather("metric"))
	want := "#[fg=colour222]\ue30d#[default] 12°C"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestWidget_Imperial(t *testing.T) {
	w := fixtureWeather("imperial")
	w.Current.TempC = 54
	got := Widget(w)
	want := "#[fg=colour222]\ue30d#[default] 54°F"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}
