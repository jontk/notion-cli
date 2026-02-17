package config

import (
	"fmt"

	"github.com/jontk/notion-cli/cmd"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration with sensitive values masked.`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()

		// Mask the API token
		maskedToken := "***"
		if len(cfg.APIToken) > 4 {
			maskedToken = cfg.APIToken[:4] + "..." + cfg.APIToken[len(cfg.APIToken)-4:]
		}

		fmt.Printf("API Token:      %s\n", maskedToken)
		fmt.Printf("Database ID:    %s\n", cfg.DatabaseID)
		fmt.Printf("Default Status: %s\n", cfg.DefaultStatus)

		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(showCmd)
}
