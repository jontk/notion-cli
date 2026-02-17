package config

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	APIToken      string `yaml:"api_token"`
	DatabaseID    string `yaml:"database_id"`
	DefaultStatus string `yaml:"default_status"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long:  `Interactive setup to create your notion-cli configuration file.`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		// Prompt for API token
		fmt.Print("Enter your Notion API token: ")
		token, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}
		token = strings.TrimSpace(token)

		// Validate token by calling user endpoint
		client := notionapi.NewClient(notionapi.Token(token))
		ctx := context.Background()
		_, err = client.User.Me(ctx)
		if err != nil {
			return fmt.Errorf("invalid API token: %w", err)
		}

		fmt.Println("✓ API token validated")

		// Prompt for database ID
		fmt.Print("Enter your Notion database ID (optional): ")
		dbID, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read database ID: %w", err)
		}
		dbID = strings.TrimSpace(dbID)

		// Prompt for default status
		fmt.Print("Enter default status for new posts [Draft]: ")
		status, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read status: %w", err)
		}
		status = strings.TrimSpace(status)
		if status == "" {
			status = "Draft"
		}

		// Create config
		cfg := Config{
			APIToken:      token,
			DatabaseID:    dbID,
			DefaultStatus: status,
		}

		// Write to config file
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		configPath := filepath.Join(home, ".notion-cli.yaml")
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		err = os.WriteFile(configPath, data, 0600)
		if err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}

		fmt.Printf("✓ Configuration saved to %s\n", configPath)
		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(initCmd)
}
