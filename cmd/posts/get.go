package posts

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
	Short: "Get a single post by ID",
	Long:  `Retrieve a single post from your Notion database by its ID.`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		cfg := cmd.GetConfig()
		client := notion.NewClient(cfg.APIToken)
		ctx := context.Background()

		if getID == "" {
			return output.Error(fmt.Errorf("post ID is required"))
		}

		post, err := client.GetPost(ctx, getID)
		if err != nil {
			return output.Error(err)
		}

		return output.JSON(post)
	},
}

func init() {
	PostsCmd.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getID, "id", "", "Post ID (required)")
	getCmd.MarkFlagRequired("id")
}
