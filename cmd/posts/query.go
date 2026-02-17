package posts

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
	queryPlatform string
	querySort     string
	queryOrder    string
	queryLimit    int
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query posts from the database",
	Long:  `Query posts from your Notion database with optional filters and sorting.`,
	Example: `  # Get all posts
  notion-cli posts query

  # Filter by status
  notion-cli posts query --status "Ready"

  # Filter by platform
  notion-cli posts query --platform "Twitter"

  # Combine filters and limit results
  notion-cli posts query --status "Draft" --platform "Blog" --limit 10

  # Sort by last edited time
  notion-cli posts query --sort "last_edited_time" --order "descending"`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		// Validate database ID
		if cfg.DatabaseID == "" {
			return output.Error(fmt.Errorf("database ID is required. Set NOTION_DATABASE_ID or run 'notion-cli config init'"))
		}

		opts := notion.QueryOptions{
			Status:   queryStatus,
			Platform: queryPlatform,
			Sort:     querySort,
			Order:    queryOrder,
			Limit:    queryLimit,
		}

		posts, err := client.QueryPosts(ctx, cfg.DatabaseID, opts)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(posts)
	},
}

func init() {
	PostsCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVar(&queryStatus, "status", "", "Filter by status (run 'databases schema' to see valid values)")
	queryCmd.Flags().StringVar(&queryPlatform, "platform", "", "Filter by platform (run 'databases schema' to see valid values)")
	queryCmd.Flags().StringVar(&querySort, "sort", "created_time", "Sort field (created_time or last_edited_time)")
	queryCmd.Flags().StringVar(&queryOrder, "order", "descending", "Sort order: ascending or descending")
	queryCmd.Flags().IntVar(&queryLimit, "limit", 100, "Maximum number of results to return")
}
