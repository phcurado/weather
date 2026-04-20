package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "weather")
	out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	return bin
}

const (
	geocodeJSON  = `{"results":[{"name":"Tallinn","country":"Estonia","latitude":59.44,"longitude":24.75,"timezone":"Europe/Tallinn"}]}`
	forecastJSON = `{
		"current":{"temperature_2m":12,"apparent_temperature":10,"relative_humidity_2m":62,"wind_speed_10m":4,"wind_direction_10m":315,"weather_code":0},
		"daily":{
			"time":["2026-04-18","2026-04-19","2026-04-20","2026-04-21","2026-04-22","2026-04-23","2026-04-24"],
			"weather_code":[0,61,2,0,0,61,2],
			"temperature_2m_max":[18,14,16,19,20,15,17],
			"temperature_2m_min":[8,9,7,8,10,9,8],
			"precipitation_sum":[0,5,1,0,0,8,2]
		}
	}`
)

func fakeServers(t *testing.T) (string, string) {
	t.Helper()
	g := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(geocodeJSON))
	}))
	t.Cleanup(g.Close)
	f := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(forecastJSON))
	}))
	t.Cleanup(f.Close)
	return g.URL, f.URL
}

func runWeather(t *testing.T, bin string, env map[string]string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Env = append([]string{}, os.Environ()...)
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			t.Fatalf("run: %v\n%s", err, stderr.String())
		}
	}
	return stdout.String(), stderr.String(), code
}

func TestCLI_Now(t *testing.T) {
	bin := buildBinary(t)
	geoBase, fcBase := fakeServers(t)

	stdout, stderr, code := runWeather(t, bin, map[string]string{
		"WEATHER_GEOCODE_BASE":  geoBase,
		"WEATHER_FORECAST_BASE": fcBase,
		"XDG_CACHE_HOME":        t.TempDir(),
		"XDG_CONFIG_HOME":       t.TempDir(),
	}, "Tallinn")
	if code != 0 {
		t.Fatalf("exit %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "Tallinn") || !strings.Contains(stdout, "Clear") {
		t.Errorf("stdout missing expected content:\n%s", stdout)
	}
}

func TestCLI_Widget_NoCityNoConfig_EmptyExit0(t *testing.T) {
	bin := buildBinary(t)
	geoBase, fcBase := fakeServers(t)

	stdout, stderr, code := runWeather(t, bin, map[string]string{
		"WEATHER_GEOCODE_BASE":  geoBase,
		"WEATHER_FORECAST_BASE": fcBase,
		"XDG_CACHE_HOME":        t.TempDir(),
		"XDG_CONFIG_HOME":       t.TempDir(),
	}, "widget")
	if code != 0 {
		t.Errorf("exit %d; want 0\nstderr: %s", code, stderr)
	}
	if stdout != "" {
		t.Errorf("stdout = %q; want empty", stdout)
	}
}

func TestCLI_Now_NoCityNoConfig_Exit1(t *testing.T) {
	bin := buildBinary(t)
	geoBase, fcBase := fakeServers(t)

	_, stderr, code := runWeather(t, bin, map[string]string{
		"WEATHER_GEOCODE_BASE":  geoBase,
		"WEATHER_FORECAST_BASE": fcBase,
		"XDG_CACHE_HOME":        t.TempDir(),
		"XDG_CONFIG_HOME":       t.TempDir(),
	})
	if code == 0 {
		t.Error("expected non-zero exit")
	}
	if !strings.Contains(stderr, "no city provided") {
		t.Errorf("stderr = %q", stderr)
	}
}
