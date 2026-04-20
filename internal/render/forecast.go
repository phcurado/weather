package render

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/phcurado/weather/internal/api"
	"github.com/phcurado/weather/internal/wmo"
)

// ANSI color escapes. Cleared to "" by DisableColor.
var (
	ansiReset  = "\033[0m"
	ansiDim    = "\033[2m"
	ansiRed    = "\033[38;5;203m"
	ansiOrange = "\033[38;5;215m"
	ansiYellow = "\033[38;5;222m"
	ansiCyan   = "\033[38;5;117m"
	ansiBlue   = "\033[38;5;111m"
	ansiGray   = "\033[38;5;244m"
)

// DisableColor strips ANSI color escapes from future render calls.
func DisableColor() {
	ansiReset = ""
	ansiDim = ""
	ansiRed = ""
	ansiOrange = ""
	ansiYellow = ""
	ansiCyan = ""
	ansiBlue = ""
	ansiGray = ""
}

// Forecast renders the 7-day table.
func Forecast(w api.Weather) string {
	t := newTable()

	header := table.Row{""}
	emojiRow := table.Row{""}
	hiRow := table.Row{colorize("High", ansiDim)}
	loRow := table.Row{colorize("Low", ansiDim)}
	rainRow := table.Row{colorize("Rain", ansiDim)}

	for _, d := range w.Daily {
		emoji, _ := wmo.Lookup(d.WeatherCode)
		header = append(header, colorize(d.Date.Weekday().String()[:3], ansiDim))
		emojiRow = append(emojiRow, emoji)
		hiRow = append(hiRow, colorize(fmt.Sprintf("%.0f°", d.TempMaxC), tempColor(d.TempMaxC, w.Units, true)))
		loRow = append(loRow, colorize(fmt.Sprintf("%.0f°", d.TempMinC), tempColor(d.TempMinC, w.Units, false)))
		rainRow = append(rainRow, formatRain(d.PrecipMM, "%.0fmm"))
	}

	t.AppendRow(header)
	t.AppendRow(emojiRow)
	t.AppendRow(hiRow)
	t.AppendRow(loRow)
	t.AppendRow(rainRow)

	t.SetColumnConfigs(columnConfigs(len(w.Daily)))
	return t.Render() + "\n"
}

// newTable builds a go-pretty table configured for this CLI's borderless style.
func newTable() table.Writer {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.Style().Options.SeparateFooter = false
	return t
}

// columnConfigs returns configs for 1 left-aligned label column + N centered data columns.
func columnConfigs(dataCols int) []table.ColumnConfig {
	cfgs := []table.ColumnConfig{{Number: 1, Align: text.AlignLeft}}
	for i := 2; i <= dataCols+1; i++ {
		cfgs = append(cfgs, table.ColumnConfig{Number: i, Align: text.AlignCenter})
	}
	return cfgs
}

// formatRain renders a colored precipitation cell, or a dim dash when dry.
func formatRain(mm float64, numFmt string) string {
	if mm <= 0 {
		return colorize("—", ansiGray)
	}
	return colorize(fmt.Sprintf(numFmt, mm), ansiCyan)
}

func colorize(s, color string) string {
	if color == "" {
		return s
	}
	return color + s + ansiReset
}

// tempColor picks an ANSI color for a temperature value. Metric bands in °C.
func tempColor(v float64, units string, isHi bool) string {
	c := v
	if units == "imperial" {
		c = (v - 32) * 5 / 9
	}
	if isHi {
		switch {
		case c >= 30:
			return ansiRed
		case c >= 20:
			return ansiOrange
		case c >= 10:
			return ansiYellow
		case c >= 0:
			return ansiCyan
		default:
			return ansiBlue
		}
	}
	switch {
	case c >= 20:
		return ansiOrange
	case c >= 10:
		return ansiYellow
	case c >= 0:
		return ansiCyan
	default:
		return ansiBlue
	}
}
