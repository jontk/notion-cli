package events

import (
	"context"
	"fmt"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "Show this week's events",
	Long:  `Show all events scheduled for the current week (Monday-Sunday).`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if cfg.EventsDatabaseID == "" {
			return output.Error(fmt.Errorf("events database ID is required"))
		}

		events, err := client.GetWeeksEvents(ctx, cfg.EventsDatabaseID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(events)
	},
}

func init() {
	EventsCmd.AddCommand(weekCmd)
}
