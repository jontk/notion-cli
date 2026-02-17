package posts

import (
	"context"
	"fmt"

	"github.com/jontk/notion-cli/cmd"
	"github.com/jontk/notion-cli/internal/notion"
	"github.com/jontk/notion-cli/internal/output"
	"github.com/spf13/cobra"
)

var archiveID string

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a post",
	Long:  `Archive a post in your Notion database.`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if archiveID == "" {
			return output.Error(fmt.Errorf("post ID is required"))
		}

		post, err := client.ArchivePost(ctx, archiveID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(post)
	},
}

func init() {
	PostsCmd.AddCommand(archiveCmd)

	archiveCmd.Flags().StringVar(&archiveID, "id", "", "Post ID (required)")
	archiveCmd.MarkFlagRequired("id")
}
