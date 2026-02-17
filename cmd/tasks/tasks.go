package tasks

import (
	"github.com/jontk/notion-cli/cmd"
	"github.com/spf13/cobra"
)

var TasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Manage tasks and TODOs in Notion",
	Long:  `Create, read, update, and complete tasks in your Notion database.`,
}

func init() {
	cmd.RootCmd.AddCommand(TasksCmd)
}
