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
	updateID        string
	updateTitle     string
	updateContent   string
	updateStatus    string
	updatePlatforms []string
	updateDate      string
	updateStdin     bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing post",
	Long:  `Update properties of an existing post in your Notion database. Only specified fields will be updated.`,
	Example: `  # Update status
  notion-cli posts update --id "PAGE_ID" --status "Published"

  # Update multiple fields
  notion-cli posts update --id "PAGE_ID" \
    --title "Updated Title" \
    --status "Ready" \
    --platform "Twitter,LinkedIn"

  # Append content
  notion-cli posts update --id "PAGE_ID" --content "Additional notes..."

  # Update from stdin
  echo '{"status":"Published"}' | notion-cli posts update --id "PAGE_ID" --stdin`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if updateID == "" {
			return output.Error(fmt.Errorf("post ID is required"))
		}

		var input models.PostInput

		if updateStdin {
			// Read JSON from stdin
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return output.Error(fmt.Errorf("failed to read stdin: %w", err))
			}
			if err := json.Unmarshal(data, &input); err != nil {
				return output.Error(fmt.Errorf("failed to parse JSON: %w", err))
			}
		} else {
			// Only set fields that were explicitly provided
			input = models.PostInput{}
			hasChanges := false

			if cobraCmd.Flags().Changed("title") {
				input.Title = updateTitle
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("content") {
				input.Content = updateContent
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("status") {
				input.Status = updateStatus
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("platform") {
				input.Platforms = updatePlatforms
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("date") {
				input.PublishDate = updateDate
				hasChanges = true
			}

			if !hasChanges {
				return output.Error(fmt.Errorf("no fields specified for update. Use --title, --status, --content, --platform, or --date to specify what to update"))
			}
		}

		post, err := client.UpdatePost(ctx, updateID, input)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(post)
	},
}

func init() {
	PostsCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateID, "id", "", "Post ID (required)")
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "New post title")
	updateCmd.Flags().StringVar(&updateContent, "content", "", "Additional content to append")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "New post status")
	updateCmd.Flags().StringSliceVar(&updatePlatforms, "platform", []string{}, "New platform(s)")
	updateCmd.Flags().StringVar(&updateDate, "date", "", "New publish date (YYYY-MM-DD)")
	updateCmd.Flags().BoolVar(&updateStdin, "stdin", false, "Read PostInput JSON from stdin")

	updateCmd.MarkFlagRequired("id")
}
