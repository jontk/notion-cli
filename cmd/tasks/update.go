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
	updateID       string
	updateTitle    string
	updateStatus   string
	updatePriority string
	updateDue      string
	updateCategory string
	updateTags     []string
	updateNotes    string
	updateStdin    bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing task",
	Long:  `Update properties of an existing task in your Notion database.`,
	Example: `  # Update status
  notion-cli tasks update --id "TASK_ID" --status "Done"

  # Update priority and due date
  notion-cli tasks update --id "TASK_ID" --priority "High" --due "2024-03-25"

  # Update from stdin
  echo '{"status":"In Progress"}' | notion-cli tasks update --id "TASK_ID" --stdin`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if updateID == "" {
			return output.Error(fmt.Errorf("task ID is required"))
		}

		var input models.TaskInput

		if updateStdin {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return output.Error(fmt.Errorf("failed to read stdin: %w", err))
			}
			if err := json.Unmarshal(data, &input); err != nil {
				return output.Error(fmt.Errorf("failed to parse JSON: %w", err))
			}
		} else {
			input = models.TaskInput{}
			hasChanges := false

			if cobraCmd.Flags().Changed("title") {
				input.Title = updateTitle
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("status") {
				input.Status = updateStatus
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("priority") {
				input.Priority = updatePriority
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("due") {
				input.DueDate = updateDue
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("category") {
				input.Category = updateCategory
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("tags") {
				input.Tags = updateTags
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

		task, err := client.UpdateTask(ctx, updateID, input)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(task)
	},
}

func init() {
	TasksCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateID, "id", "", "Task ID (required)")
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "New task title")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "New status")
	updateCmd.Flags().StringVar(&updatePriority, "priority", "", "New priority")
	updateCmd.Flags().StringVar(&updateDue, "due", "", "New due date (YYYY-MM-DD)")
	updateCmd.Flags().StringVar(&updateCategory, "category", "", "New category")
	updateCmd.Flags().StringSliceVar(&updateTags, "tags", []string{}, "New tags")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "New notes")
	updateCmd.Flags().BoolVar(&updateStdin, "stdin", false, "Read TaskInput JSON from stdin")

	updateCmd.MarkFlagRequired("id")
}
