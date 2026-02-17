package tasks

import (
	"context"
	"fmt"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	queryStatus   string
	queryPriority string
	queryCategory string
	queryLimit    int
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query tasks from the database",
	Long:  `Query tasks from your Notion database with optional filters.`,
	Example: `  # All tasks
  notion-cli tasks query

  # High priority tasks
  notion-cli tasks query --priority "High"

  # Work tasks
  notion-cli tasks query --category "Work"

  # Todo items
  notion-cli tasks query --status "Todo" --limit 10`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if cfg.TasksDatabaseID == "" {
			return output.Error(fmt.Errorf("tasks database ID is required"))
		}

		opts := notion.TaskQueryOptions{
			Status:   queryStatus,
			Priority: queryPriority,
			Category: queryCategory,
			Limit:    queryLimit,
		}

		tasks, err := client.QueryTasks(ctx, cfg.TasksDatabaseID, opts)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(tasks)
	},
}

func init() {
	TasksCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVar(&queryStatus, "status", "", "Filter by status (Todo, In Progress, Done, Blocked)")
	queryCmd.Flags().StringVar(&queryPriority, "priority", "", "Filter by priority (High, Medium, Low)")
	queryCmd.Flags().StringVar(&queryCategory, "category", "", "Filter by category")
	queryCmd.Flags().IntVar(&queryLimit, "limit", 100, "Maximum number of results")
}
