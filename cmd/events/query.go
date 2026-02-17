package events

import (
	"context"
	"fmt"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	queryType   string
	queryStatus string
	queryLimit  int
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query events from the database",
	Long:  `Query events from your Notion calendar database with optional filters.`,
	Example: `  # All events
  notion-cli events query

  # Work events
  notion-cli events query --type "Work"

  # Scheduled events
  notion-cli events query --status "Scheduled"

  # Limited results
  notion-cli events query --limit 20`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if cfg.EventsDatabaseID == "" {
			return output.Error(fmt.Errorf("events database ID is required"))
		}

		opts := notion.EventQueryOptions{
			Type:   queryType,
			Status: queryStatus,
			Limit:  queryLimit,
		}

		events, err := client.QueryEvents(ctx, cfg.EventsDatabaseID, opts)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(events)
	},
}

func init() {
	EventsCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVar(&queryType, "type", "", "Filter by event type")
	queryCmd.Flags().StringVar(&queryStatus, "status", "", "Filter by status (Scheduled, Completed, Cancelled)")
	queryCmd.Flags().IntVar(&queryLimit, "limit", 100, "Maximum number of results")
}
