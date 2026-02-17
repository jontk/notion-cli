package events

import (
	"context"
	"fmt"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var getID string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a single event by ID",
	Long:  `Retrieve a single event from your Notion calendar database by its ID.`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if getID == "" {
			return output.Error(fmt.Errorf("event ID is required"))
		}

		event, err := client.GetEvent(ctx, getID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(event)
	},
}

func init() {
	EventsCmd.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getID, "id", "", "Event ID (required)")
	getCmd.MarkFlagRequired("id")
}
