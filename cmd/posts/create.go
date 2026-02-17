package posts

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
	createTitle         string
	createContent       string
	createStatus        string
	createWeek          int
	createPillar        string
	createDueDate       string
	createDistributedTo []string
	createHashtags      []string
	createStdin         bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new post",
	Long:  `Create a new post in your Notion content pipeline database.`,
	Example: `  # Create a draft post
  notion-cli posts create --title "Introducing s9s: a terminal UI for SLURM"

  # Create with full details
  notion-cli posts create \
    --title "Introducing s9s" \
    --status "Outline" \
    --pillar "SLURM & HPC" \
    --week 1 \
    --due-date "2026-02-24"

  # Create from stdin (AI workflow)
  echo '{"title":"AI Post","status":"Draft","pillar":"Go Tools","week":2}' | \
    notion-cli posts create --stdin`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		var input models.PostInput

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
			input = models.PostInput{
				Title:         createTitle,
				Content:       createContent,
				Status:        createStatus,
				Week:          createWeek,
				Pillar:        createPillar,
				PublishDate:   createDueDate,
				DistributedTo: createDistributedTo,
				Hashtags:      createHashtags,
			}
		}

		if input.Status == "" {
			input.Status = cfg.DefaultStatus
		}

		if cfg.DatabaseID == "" {
			return output.Error(fmt.Errorf("database ID is required. Set NOTION_DATABASE_ID or run 'notion-cli config init'"))
		}

		post, err := client.CreatePost(ctx, input, cfg.DatabaseID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(post)
	},
}

func init() {
	PostsCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createTitle, "title", "", "Post title (required unless using --stdin)")
	createCmd.Flags().StringVar(&createContent, "content", "", "Post content (markdown supported)")
	createCmd.Flags().StringVar(&createStatus, "status", "", "Status: Idea, Outline, Draft, Review, Published, Distributed (default: Draft)")
	createCmd.Flags().IntVar(&createWeek, "week", 0, "Week number in the content calendar (1-12)")
	createCmd.Flags().StringVar(&createPillar, "pillar", "", "Content pillar: 'SLURM & HPC', 'Go Tools', 'Infrastructure', 'Career & AI'")
	createCmd.Flags().StringVar(&createDueDate, "due-date", "", "Target publish date in YYYY-MM-DD format")
	createCmd.Flags().StringSliceVar(&createDistributedTo, "distributed-to", []string{}, "Platforms distributed to: LinkedIn,Twitter,Dev.to,Hacker News,Reddit")
	createCmd.Flags().StringSliceVar(&createHashtags, "hashtags", []string{}, "Hashtags (comma-separated)")
	createCmd.Flags().BoolVar(&createStdin, "stdin", false, "Read PostInput JSON from stdin (for AI/automation workflows)")
}
