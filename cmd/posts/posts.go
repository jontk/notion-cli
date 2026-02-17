package posts

import (
	"github.com/jontk/notion-cli/cmd"
	"github.com/spf13/cobra"
)

var PostsCmd = &cobra.Command{
	Use:   "posts",
	Short: "Manage posts in Notion",
	Long:  `Create, read, update, and archive posts in your Notion database.`,
}

func init() {
	cmd.RootCmd.AddCommand(PostsCmd)
}
