# weather

Tiny Go CLI for current conditions, a 7-day forecast, hourly breakdown, and a tmux statusbar widget. Backed by [Open-Meteo](https://open-meteo.com). No API key.

## Install

```
curl -sSL https://raw.githubusercontent.com/phcurado/weather/main/install.sh | sh
```

Or: `go install github.com/phcurado/weather/cmd/weather@latest`.

## Usage

```
weather [city]                # summary + 7-day forecast
weather [city] --hourly       # summary + next 12 hours
weather [city] -n 24          # summary + N hours (implies hourly)
weather widget                # tmux status-line line
weather config                # print resolved config + path
```

If `[city]` is omitted, `default_city` from config is used.

Flags:
- `-H, --hourly` — show hourly view instead of the 7-day table
- `-n, --hours N` — hours to show with hourly view (default 12)

## Configuration

Path: `$XDG_CONFIG_HOME/weather/config.toml` (fallback `~/.config/weather/config.toml`).

```toml
default_city = "Tallinn"
units        = "metric"   # "metric" | "imperial"
cache_ttl    = "10m"
```

All fields optional. Defaults: no default city, metric, 10m cache TTL.

## tmux

```
set -g status-right '#(command -v weather >/dev/null 2>&1 && weather widget || echo "")  %I:%M %p  %a %b %d'
```

`widget` is designed to fail silently — the segment renders empty if the binary is missing or the network is down. It emits tmux format strings (`#[fg=...]`) so tmux colorizes the glyph natively. Requires a Nerd Font patched terminal font.

## Cache

- `$XDG_CACHE_HOME/weather/geocode/` — permanent
- `$XDG_CACHE_HOME/weather/weather/` — TTL from config
- On API failure with any cache hit (fresh or stale): the cache is used; interactive commands emit a warning on stderr, `widget` stays silent.

## License

MIT.
