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
	queryStatus        string
	queryPillar        string
	queryDistributedTo string
	querySort          string
	queryOrder         string
	queryLimit         int
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query posts from the content pipeline",
	Long:  `Query posts from your Notion content pipeline with optional filters and sorting.`,
	Example: `  # Get all posts
  notion-cli posts query

  # Filter by status
  notion-cli posts query --status "Draft"
  notion-cli posts query --status "Review"

  # Filter by content pillar
  notion-cli posts query --pillar "SLURM & HPC"
  notion-cli posts query --pillar "Go Tools"

  # Filter by distribution platform
  notion-cli posts query --distributed-to "LinkedIn"

  # Combine filters
  notion-cli posts query --status "Published" --pillar "Infrastructure"

  # Sort by last edited
  notion-cli posts query --sort "last_edited_time" --order "descending"`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if cfg.DatabaseID == "" {
			return output.Error(fmt.Errorf("database ID is required. Set NOTION_DATABASE_ID or run 'notion-cli config init'"))
		}

		opts := notion.QueryOptions{
			Status:        queryStatus,
			Pillar:        queryPillar,
			DistributedTo: queryDistributedTo,
			Sort:          querySort,
			Order:         queryOrder,
			Limit:         queryLimit,
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

	queryCmd.Flags().StringVar(&queryStatus, "status", "", "Filter by status: Idea, Outline, Draft, Review, Published, Distributed")
	queryCmd.Flags().StringVar(&queryPillar, "pillar", "", "Filter by pillar: 'SLURM & HPC', 'Go Tools', 'Infrastructure', 'Career & AI'")
	queryCmd.Flags().StringVar(&queryDistributedTo, "distributed-to", "", "Filter by distribution platform: LinkedIn, Twitter, Dev.to, Hacker News, Reddit")
	queryCmd.Flags().StringVar(&querySort, "sort", "created_time", "Sort field: created_time or last_edited_time")
	queryCmd.Flags().StringVar(&queryOrder, "order", "descending", "Sort order: ascending or descending")
	queryCmd.Flags().IntVar(&queryLimit, "limit", 100, "Maximum number of results")
}
