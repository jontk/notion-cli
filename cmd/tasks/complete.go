package tasks

import (
	"context"
	"fmt"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var completeID string

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark a task as complete",
	Long:  `Mark a task as complete (Done status) in your Notion database.`,
	Example: `  notion-cli tasks complete --id "TASK_ID"`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if completeID == "" {
			return output.Error(fmt.Errorf("task ID is required"))
		}

		task, err := client.CompleteTask(ctx, completeID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(task)
	},
}

func init() {
	TasksCmd.AddCommand(completeCmd)

	completeCmd.Flags().StringVar(&completeID, "id", "", "Task ID (required)")
	completeCmd.MarkFlagRequired("id")
}
