package main

import (
	"fmt"

	"github.com/phcurado/weather/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Print resolved config and config file path",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			fmt.Printf("path      %s\n", config.Path())
			fmt.Printf("city      %s\n", cfg.City)
			fmt.Printf("units     %s\n", cfg.Units)
			fmt.Printf("cache_ttl %s\n", cfg.CacheTTL)
			return nil
		},
	}
}
