package main

import (
	"os"

	"github.com/jontk/notion-cli/cmd"
	_ "github.com/jontk/notion-cli/cmd/config"
	_ "github.com/jontk/notion-cli/cmd/databases"
	_ "github.com/jontk/notion-cli/cmd/posts"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
