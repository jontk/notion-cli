package tasks

import (
	"context"
	"fmt"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var overdueCmd = &cobra.Command{
	Use:   "overdue",
	Short: "Show overdue tasks",
	Long:  `Show all tasks that are past their due date and still incomplete.`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if cfg.TasksDatabaseID == "" {
			return output.Error(fmt.Errorf("tasks database ID is required"))
		}

		tasks, err := client.GetOverdueTasks(ctx, cfg.TasksDatabaseID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(tasks)
	},
}

func init() {
	TasksCmd.AddCommand(overdueCmd)
}
