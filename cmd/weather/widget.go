package main

import (
	"fmt"
	"io"
	"os"

	"github.com/phcurado/weather/internal/render"
	"github.com/spf13/cobra"
)

func newWidgetCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "widget [city]",
		Short:         "Single-line tmux widget",
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(_ *cobra.Command, args []string) error {
			city := ""
			if len(args) == 1 {
				city = args[0]
			}
			r, err := resolve(city, io.Discard)
			if err != nil {
				os.Exit(0)
			}
			fmt.Print(render.Widget(r.Weather))
			return nil
		},
	}
}
