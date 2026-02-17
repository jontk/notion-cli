package events

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/models"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	updateID        string
	updateTitle     string
	updateDate      string
	updateType      string
	updateLocation  string
	updateAttendees []string
	updateStatus    string
	updateNotes     string
	updateStdin     bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing event",
	Long:  `Update properties of an existing event in your Notion calendar database.`,
	Example: `  # Update status
  notion-cli events update --id "EVENT_ID" --status "Completed"

  # Update date and location
  notion-cli events update --id "EVENT_ID" --date "2024-03-26 15:00" --location "Room B"

  # Update from stdin
  echo '{"status":"Cancelled"}' | notion-cli events update --id "EVENT_ID" --stdin`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if updateID == "" {
			return output.Error(fmt.Errorf("event ID is required"))
		}

		var input models.EventInput

		if updateStdin {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return output.Error(fmt.Errorf("failed to read stdin: %w", err))
			}
			if err := json.Unmarshal(data, &input); err != nil {
				return output.Error(fmt.Errorf("failed to parse JSON: %w", err))
			}
		} else {
			input = models.EventInput{}
			hasChanges := false

			if cobraCmd.Flags().Changed("title") {
				input.Title = updateTitle
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("date") {
				input.Date = updateDate
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("type") {
				input.Type = updateType
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("location") {
				input.Location = updateLocation
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("attendees") {
				input.Attendees = updateAttendees
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("status") {
				input.Status = updateStatus
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("notes") {
				input.Notes = updateNotes
				hasChanges = true
			}

			if !hasChanges {
				return output.Error(fmt.Errorf("no fields specified for update"))
			}
		}

		event, err := client.UpdateEvent(ctx, updateID, input)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(event)
	},
}

func init() {
	EventsCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateID, "id", "", "Event ID (required)")
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "New event title")
	updateCmd.Flags().StringVar(&updateDate, "date", "", "New event date and time (YYYY-MM-DD HH:MM)")
	updateCmd.Flags().StringVar(&updateType, "type", "", "New event type")
	updateCmd.Flags().StringVar(&updateLocation, "location", "", "New location")
	updateCmd.Flags().StringSliceVar(&updateAttendees, "attendees", []string{}, "New attendees")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "New status")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "New notes")
	updateCmd.Flags().BoolVar(&updateStdin, "stdin", false, "Read EventInput JSON from stdin")

	updateCmd.MarkFlagRequired("id")
}
