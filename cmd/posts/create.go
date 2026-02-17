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
	createTitle       string
	createContent     string
	createStatus      string
	createPlatforms   []string
	createDate        string
	createStdin       bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new post",
	Long: `Create a new post in your Notion database with the specified properties.`,
	Example: `  # Create a simple draft post
  notion-cli posts create --title "My First Post"

  # Create with all fields
  notion-cli posts create \
    --title "Blog Post About Go" \
    --content "Go is a great language for CLI tools..." \
    --status "Ready" \
    --platform "Twitter,LinkedIn,Blog" \
    --date "2024-03-20"

  # Create from stdin (useful for AI-generated content)
  echo '{"title":"AI Post","content":"Generated content..."}' | \
    notion-cli posts create --stdin`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		var input models.PostInput

		if createStdin {
			// Read JSON from stdin
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return output.Error(fmt.Errorf("failed to read stdin: %w", err))
			}
			if err := json.Unmarshal(data, &input); err != nil {
				return output.Error(fmt.Errorf("failed to parse JSON: %w", err))
			}
		} else {
			// Use flags
			if createTitle == "" {
				return output.Error(fmt.Errorf("title is required"))
			}

			input = models.PostInput{
				Title:       createTitle,
				Content:     createContent,
				Status:      createStatus,
				Platforms:   createPlatforms,
				PublishDate: createDate,
			}
		}

		// Set default status if not provided
		if input.Status == "" {
			input.Status = cfg.DefaultStatus
		}

		// Validate database ID
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
	createCmd.Flags().StringVar(&createStatus, "status", "", "Post status - use 'databases schema' to see valid values (default: Draft)")
	createCmd.Flags().StringSliceVar(&createPlatforms, "platform", []string{}, "Platform(s) - comma-separated (e.g., Twitter,LinkedIn,Blog)")
	createCmd.Flags().StringVar(&createDate, "date", "", "Publish date in YYYY-MM-DD format (e.g., 2024-03-20)")
	createCmd.Flags().BoolVar(&createStdin, "stdin", false, "Read PostInput JSON from stdin (for AI/automation workflows)")
}
