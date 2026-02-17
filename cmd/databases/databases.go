package databases

import (
	"github.com/jontk/notion-cli/cmd"
	"github.com/spf13/cobra"
)

var DatabasesCmd = &cobra.Command{
	Use:   "databases",
	Short: "Manage databases in Notion",
	Long:  `List and inspect databases in your Notion workspace.`,
}

func init() {
	cmd.RootCmd.AddCommand(DatabasesCmd)
}
