package tasks

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
	createTitle    string
	createStatus   string
	createPriority string
	createDue      string
	createCategory string
	createTags     []string
	createNotes    string
	createStdin    bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Long:  `Create a new task in your Notion database.`,
	Example: `  # Simple task
  notion-cli tasks create --title "Buy milk"

  # With all fields
  notion-cli tasks create \
    --title "Review PR #123" \
    --category "Work" \
    --priority "High" \
    --due "2024-03-20" \
    --tags "urgent,review"

  # From stdin
  echo '{"title":"Call dentist","category":"Personal"}' | \
    notion-cli tasks create --stdin`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		var input models.TaskInput

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

			input = models.TaskInput{
				Title:    createTitle,
				Status:   createStatus,
				Priority: createPriority,
				DueDate:  createDue,
				Category: createCategory,
				Tags:     createTags,
				Notes:    createNotes,
			}
		}

		// Set defaults
		if input.Status == "" {
			input.Status = cfg.DefaultTaskStatus
		}
		if input.Priority == "" {
			input.Priority = cfg.DefaultPriority
		}

		if cfg.TasksDatabaseID == "" {
			return output.Error(fmt.Errorf("tasks database ID is required. Set NOTION_TASKS_DATABASE_ID or add to config"))
		}

		task, err := client.CreateTask(ctx, input, cfg.TasksDatabaseID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(task)
	},
}

func init() {
	TasksCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createTitle, "title", "", "Task title (required unless using --stdin)")
	createCmd.Flags().StringVar(&createStatus, "status", "", "Task status (default: Todo)")
	createCmd.Flags().StringVar(&createPriority, "priority", "", "Priority: High, Medium, Low (default: Medium)")
	createCmd.Flags().StringVar(&createDue, "due", "", "Due date in YYYY-MM-DD format or 'today', 'tomorrow'")
	createCmd.Flags().StringVar(&createCategory, "category", "", "Category: Work, Personal, Shopping, Project, etc.")
	createCmd.Flags().StringSliceVar(&createTags, "tags", []string{}, "Tags (comma-separated)")
	createCmd.Flags().StringVar(&createNotes, "notes", "", "Additional notes")
	createCmd.Flags().BoolVar(&createStdin, "stdin", false, "Read TaskInput JSON from stdin")
}
