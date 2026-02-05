package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/pkg/tui"
)

func main() {
	config := tui.TUIConfig{
		StatisticsStoragePath: "./.pm/stats.json",
		BeadsDBPath:           "./.pm/db.db",
		IssuePrefix:           "pm",
	}

	if err := tui.Run(context.Background(), config); err != nil {
		panic(err)
	}
}
