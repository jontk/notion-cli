package databases

import (
	"context"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all databases",
	Long:  `List all databases accessible to your Notion integration.`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		databases, err := client.ListDatabases(ctx)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(databases)
	},
}

func init() {
	DatabasesCmd.AddCommand(listCmd)
}
