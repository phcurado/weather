package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/phcurado/weather/internal/render"
)

// version is set via -ldflags at release time.
var version = "dev"

func main() {
	if !stdoutIsTTY() || os.Getenv("NO_COLOR") != "" {
		render.DisableColor()
	}

	var (
		hourly bool
		hours  int
	)

	root := &cobra.Command{
		Use:     "weather [city]",
		Short:   "Tiny weather CLI backed by Open-Meteo",
		Version: version,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			city := ""
			if len(args) == 1 {
				city = args[0]
			}
			r, err := resolve(city, os.Stderr)
			if err != nil {
				return err
			}
			fmt.Print(render.Summary(r.Weather))
			fmt.Println()
			if hourly || cmd.Flags().Changed("hours") {
				fmt.Print(render.Hourly(r.Weather, hours))
			} else {
				fmt.Print(render.Forecast(r.Weather))
			}
			return nil
		},
	}
	root.Flags().BoolVarP(&hourly, "hourly", "H", false, "show next-N-hours view instead of 7-day")
	root.Flags().IntVarP(&hours, "hours", "n", 12, "hours to show with --hourly")

	root.AddCommand(newConfigCmd())
	root.AddCommand(newWidgetCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
