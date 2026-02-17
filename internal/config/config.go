package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	APIToken      string
	DatabaseID    string
	DefaultStatus string
}

func Load() (*Config, error) {
	cfg := &Config{
		APIToken:      viper.GetString("api_token"),
		DatabaseID:    viper.GetString("database_id"),
		DefaultStatus: viper.GetString("default_status"),
	}

	// Set default status if not configured
	if cfg.DefaultStatus == "" {
		cfg.DefaultStatus = "Draft"
	}

	// Validate required fields
	if cfg.APIToken == "" {
		return nil, fmt.Errorf("API token is required. Set NOTION_API_TOKEN environment variable or run 'notion-cli config init'")
	}

	return cfg, nil
}
