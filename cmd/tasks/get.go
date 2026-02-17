package tasks

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
	Short: "Get a single task by ID",
	Long:  `Retrieve a single task from your Notion database by its ID.`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if getID == "" {
			return output.Error(fmt.Errorf("task ID is required"))
		}

		task, err := client.GetTask(ctx, getID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(task)
	},
}

func init() {
	TasksCmd.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getID, "id", "", "Task ID (required)")
	getCmd.MarkFlagRequired("id")
}
