package events

import (
	"github.com/jontk/notion-cli/cmd"
	"github.com/spf13/cobra"
)

var EventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Manage calendar events",
	Long:  `Manage calendar events in your Notion database. Create, query, update, and organize your schedule.`,
}

func init() {
	cmd.RootCmd.AddCommand(EventsCmd)
}
