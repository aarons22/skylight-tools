package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skylight",
	Short: "Skylight Calendar API CLI",
	Long:  "A command-line interface for the Skylight Calendar API.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "Output raw JSON")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().String("base-url", "https://api.ourskylight.com/api", "API base URL")
	rootCmd.PersistentFlags().String("frame-id", "", "Skylight frame ID (overrides SKYLIGHT_FRAME_ID env var and config)")
	rootCmd.PersistentFlags().String("config", "", "Config file path (default: ~/.config/skylight/config.yaml)")
}
