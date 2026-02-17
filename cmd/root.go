package cmd

import (
	"fmt"
	"os"

	"github.com/jontk/notion-cli/internal/config"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	outputFormat string
	cfg          *config.Config
	version      = "0.2.0"
)

var RootCmd = &cobra.Command{
	Use:     "notion-cli",
	Short:   "A CLI tool for managing Notion content, tasks, and events",
	Version: version,
	Long: `notion-cli is a command-line interface for managing content, tasks, and events in Notion databases.
It provides an easy way to create, read, update, and organize posts, tasks, and calendar events.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(); err != nil {
			return err
		}
		return nil
	},
}

var rootCmd = RootCmd

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.notion-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "output", "json", "output format (json|table)")
}

func initConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".notion-cli")
	}

	viper.SetEnvPrefix("NOTION")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var err error
	cfg, err = config.Load()
	if err != nil {
		return output.Error(err)
	}

	return nil
}

func GetConfig() *config.Config {
	return cfg
}

func GetOutputFormat() string {
	return outputFormat
}
