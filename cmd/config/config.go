package config

import (
	"github.com/jontk/notion-cli/cmd"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Initialize and manage notion-cli configuration.`,
}

func init() {
	cmd.RootCmd.AddCommand(ConfigCmd)
}
