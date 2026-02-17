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
	updateID            string
	updateTitle         string
	updateContent       string
	updateStatus        string
	updateWeek          int
	updatePillar        string
	updatePublishDate   string
	updatePublishedDate string
	updateBlogURL       string
	updateDistributedTo []string
	updateDistributedDate string
	updateLinkedInDraft string
	updateTwitterThread string
	updateHNTitle       string
	updateRedditTitle   string
	updateHashtags      []string
	updateStdin         bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing post",
	Long:  `Update properties of an existing post. Only specified fields will be updated.`,
	Example: `  # Advance status
  notion-cli posts update --id "PAGE_ID" --status "Draft"

  # Mark as published
  notion-cli posts update --id "PAGE_ID" \
    --status "Published" \
    --published-date "2026-02-24" \
    --blog-url "https://jontk.com/blog/introducing-s9s"

  # Mark as distributed
  notion-cli posts update --id "PAGE_ID" \
    --status "Distributed" \
    --distributed-to "LinkedIn,Twitter,Dev.to" \
    --distributed-date "2026-02-24"

  # Update from stdin
  echo '{"status":"Review"}' | notion-cli posts update --id "PAGE_ID" --stdin`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if updateID == "" {
			return output.Error(fmt.Errorf("post ID is required"))
		}

		var input models.PostInput

		if updateStdin {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return output.Error(fmt.Errorf("failed to read stdin: %w", err))
			}
			if err := json.Unmarshal(data, &input); err != nil {
				return output.Error(fmt.Errorf("failed to parse JSON: %w", err))
			}
		} else {
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
			if cobraCmd.Flags().Changed("week") {
				input.Week = updateWeek
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("pillar") {
				input.Pillar = updatePillar
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("publish-date") {
				input.PublishDate = updatePublishDate
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("published-date") {
				input.PublishedDate = updatePublishedDate
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("blog-url") {
				input.BlogURL = updateBlogURL
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("distributed-to") {
				input.DistributedTo = updateDistributedTo
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("distributed-date") {
				input.DistributedDate = updateDistributedDate
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("linkedin-draft") {
				input.LinkedInDraft = updateLinkedInDraft
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("twitter-thread") {
				input.TwitterThread = updateTwitterThread
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("hn-title") {
				input.HNTitle = updateHNTitle
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("reddit-title") {
				input.RedditTitle = updateRedditTitle
				hasChanges = true
			}
			if cobraCmd.Flags().Changed("hashtags") {
				input.Hashtags = updateHashtags
				hasChanges = true
			}

			if !hasChanges {
				return output.Error(fmt.Errorf("no fields specified for update"))
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
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "New title")
	updateCmd.Flags().StringVar(&updateContent, "content", "", "Content to append")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "Status: Idea, Outline, Draft, Review, Published, Distributed")
	updateCmd.Flags().IntVar(&updateWeek, "week", 0, "Week number in the content calendar")
	updateCmd.Flags().StringVar(&updatePillar, "pillar", "", "Content pillar")
	updateCmd.Flags().StringVar(&updatePublishDate, "publish-date", "", "Target publish date (YYYY-MM-DD)")
	updateCmd.Flags().StringVar(&updatePublishedDate, "published-date", "", "Actual publish date (YYYY-MM-DD)")
	updateCmd.Flags().StringVar(&updateBlogURL, "blog-url", "", "Blog post URL")
	updateCmd.Flags().StringSliceVar(&updateDistributedTo, "distributed-to", []string{}, "Platforms distributed to (comma-separated)")
	updateCmd.Flags().StringVar(&updateDistributedDate, "distributed-date", "", "Distribution date (YYYY-MM-DD)")
	updateCmd.Flags().StringVar(&updateLinkedInDraft, "linkedin-draft", "", "LinkedIn post draft")
	updateCmd.Flags().StringVar(&updateTwitterThread, "twitter-thread", "", "Twitter thread draft")
	updateCmd.Flags().StringVar(&updateHNTitle, "hn-title", "", "Hacker News title")
	updateCmd.Flags().StringVar(&updateRedditTitle, "reddit-title", "", "Reddit title")
	updateCmd.Flags().StringSliceVar(&updateHashtags, "hashtags", []string{}, "Hashtags (comma-separated)")
	updateCmd.Flags().BoolVar(&updateStdin, "stdin", false, "Read PostInput JSON from stdin")

	updateCmd.MarkFlagRequired("id")
}
