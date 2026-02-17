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
	createTitle     string
	createDate      string
	createType      string
	createLocation  string
	createAttendees []string
	createStatus    string
	createNotes     string
	createStdin     bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new event",
	Long:  `Create a new event in your Notion calendar database.`,
	Example: `  # Simple event
  notion-cli events create --title "Team Meeting" --date "2024-03-20 14:00"

  # Full event details
  notion-cli events create \
    --title "Product Launch" \
    --date "2024-04-15 10:00" \
    --type "Work" \
    --location "Conference Room A" \
    --attendees "john@example.com,jane@example.com" \
    --status "Scheduled" \
    --notes "Q2 product release event"

  # From stdin
  echo '{"title":"Doctor Appointment","date":"2024-03-25 09:30","type":"Personal"}' | \
    notion-cli events create --stdin`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if cfg.EventsDatabaseID == "" {
			return output.Error(fmt.Errorf("events database ID is required"))
		}

		var input models.EventInput

		if createStdin {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return output.Error(fmt.Errorf("failed to read stdin: %w", err))
			}
			if err := json.Unmarshal(data, &input); err != nil {
				return output.Error(fmt.Errorf("failed to parse JSON: %w", err))
			}
		} else {
			if createTitle == "" {
				return output.Error(fmt.Errorf("title is required"))
			}
			if createDate == "" {
				return output.Error(fmt.Errorf("date is required"))
			}

			input = models.EventInput{
				Title:     createTitle,
				Date:      createDate,
				Type:      createType,
				Location:  createLocation,
				Attendees: createAttendees,
				Status:    createStatus,
				Notes:     createNotes,
			}
		}

		event, err := client.CreateEvent(ctx, input, cfg.EventsDatabaseID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(event)
	},
}

func init() {
	EventsCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createTitle, "title", "", "Event title (required)")
	createCmd.Flags().StringVar(&createDate, "date", "", "Event date and time (YYYY-MM-DD HH:MM)")
	createCmd.Flags().StringVar(&createType, "type", "", "Event type")
	createCmd.Flags().StringVar(&createLocation, "location", "", "Event location")
	createCmd.Flags().StringSliceVar(&createAttendees, "attendees", []string{}, "Event attendees (comma-separated)")
	createCmd.Flags().StringVar(&createStatus, "status", "", "Event status")
	createCmd.Flags().StringVar(&createNotes, "notes", "", "Event notes")
	createCmd.Flags().BoolVar(&createStdin, "stdin", false, "Read EventInput JSON from stdin")
}
