package events

import (
	"context"
	"fmt"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var cancelID string

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel an event",
	Long:  `Mark an event as cancelled in your Notion calendar database.`,
	Example: `  notion-cli events cancel --id "EVENT_ID"`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if cancelID == "" {
			return output.Error(fmt.Errorf("event ID is required"))
		}

		event, err := client.CancelEvent(ctx, cancelID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(event)
	},
}

func init() {
	EventsCmd.AddCommand(cancelCmd)

	cancelCmd.Flags().StringVar(&cancelID, "id", "", "Event ID (required)")
	cancelCmd.MarkFlagRequired("id")
}
