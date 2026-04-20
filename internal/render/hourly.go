package render

import (
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/phcurado/weather/internal/api"
	"github.com/phcurado/weather/internal/wmo"
)

// Hourly renders the next `hours` hours from now as a horizontal table.
func Hourly(w api.Weather, hours int) string {
	picks := pickHours(w.Hourly, hours)
	if len(picks) == 0 {
		return ""
	}

	t := newTable()

	header := table.Row{""}
	emojis := table.Row{""}
	temps := table.Row{colorize("Temp", ansiDim)}
	rains := table.Row{colorize("Rain", ansiDim)}

	for _, h := range picks {
		emoji, _ := wmo.Lookup(h.WeatherCode)
		header = append(header, colorize(h.Time.Format("15"), ansiDim))
		emojis = append(emojis, emoji)
		temps = append(temps, colorize(fmt.Sprintf("%.0f°", h.TempC), tempColor(h.TempC, w.Units, true)))
		rains = append(rains, formatRain(h.PrecipMM, "%.1fmm"))
	}

	t.AppendRow(header)
	t.AppendRow(emojis)
	t.AppendRow(temps)
	t.AppendRow(rains)

	t.SetColumnConfigs(columnConfigs(len(picks)))
	return t.Render() + "\n"
}

// pickHours returns up to `count` entries from now onward (skipping past hours).
func pickHours(all []api.Hourly, count int) []api.Hourly {
	start := 0
	cutoff := time.Now().Truncate(time.Hour)
	for i, h := range all {
		if !h.Time.Before(cutoff) {
			start = i
			break
		}
	}
	end := start + count
	if end > len(all) {
		end = len(all)
	}
	return all[start:end]
}
