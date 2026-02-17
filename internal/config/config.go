package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	APIToken           string
	DatabaseID         string
	TasksDatabaseID    string
	EventsDatabaseID   string
	DefaultStatus      string
	DefaultTaskStatus  string
	DefaultPriority    string
}

func Load() (*Config, error) {
	cfg := &Config{
		APIToken:          viper.GetString("api_token"),
		DatabaseID:        viper.GetString("database_id"),
		TasksDatabaseID:   viper.GetString("tasks_database_id"),
		EventsDatabaseID:  viper.GetString("events_database_id"),
		DefaultStatus:     viper.GetString("default_status"),
		DefaultTaskStatus: viper.GetString("default_task_status"),
		DefaultPriority:   viper.GetString("default_priority"),
	}

	// Set defaults if not configured
	if cfg.DefaultStatus == "" {
		cfg.DefaultStatus = "Draft"
	}
	if cfg.DefaultTaskStatus == "" {
		cfg.DefaultTaskStatus = "Todo"
	}
	if cfg.DefaultPriority == "" {
		cfg.DefaultPriority = "Medium"
	}

	// Validate required fields
	if cfg.APIToken == "" {
		return nil, fmt.Errorf("API token is required. Set NOTION_API_TOKEN environment variable or run 'notion-cli config init'")
	}

	return cfg, nil
}
