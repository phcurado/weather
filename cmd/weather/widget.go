package main

import (
	"fmt"
	"io"
	"os"

	"github.com/phcurado/weather/internal/render"
	"github.com/spf13/cobra"
)

func newWidgetCmd() *cobra.Command {
	var here bool
	cmd := &cobra.Command{
		Use:           "widget [city]",
		Short:         "Single-line tmux widget",
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			city := ""
			if len(args) == 1 {
				city = args[0]
			}
			if here && city != "" {
				os.Exit(0)
			}
			r, err := resolve(city, here, io.Discard)
			if err != nil {
				os.Exit(0)
			}
			fmt.Print(render.Widget(r.Weather))
			return nil
		},
	}
	cmd.Flags().BoolVarP(&here, "here", "l", false, "use current location from IP geolocation")
	return cmd
}
