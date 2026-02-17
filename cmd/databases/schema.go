package databases

import (
	"context"
	"fmt"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var schemaID string

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Get database schema",
	Long:  `Retrieve the schema (properties and their types) of a Notion database.`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		// Use config database ID if not provided
		databaseID := schemaID
		if databaseID == "" {
			databaseID = cfg.DatabaseID
		}

		if databaseID == "" {
			return output.Error(fmt.Errorf("database ID is required"))
		}

		schema, err := client.GetSchema(ctx, databaseID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(schema)
	},
}

func init() {
	DatabasesCmd.AddCommand(schemaCmd)

	schemaCmd.Flags().StringVar(&schemaID, "id", "", "Database ID (defaults to configured database)")
}
