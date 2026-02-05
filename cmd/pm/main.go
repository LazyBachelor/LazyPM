package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/pkg/cli"
)

func main() {
	config := cli.CLIConfig{
		RootCmd:               "pm",
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/stats.json",
	}

	if err := cli.Run(context.Background(), config); err != nil {
		return
	}
}
