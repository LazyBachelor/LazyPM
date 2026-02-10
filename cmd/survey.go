package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LazyBachelor/LazyPM/pkg"
	"github.com/LazyBachelor/LazyPM/pkg/cli"
	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
)

func main() {
	config := pkg.SurveyConfig{
		RootCmd:               "pm",
		WebAddress:            "localhost:8080",
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/stats.json",
	}

	ctx := context.Background()
	var err error

	switch os.Args[1] {
	case "tui":
		_, err = tui.Run(ctx, config)
	case "cli":
		err = cli.RunWithArgs(ctx, config, os.Args[2:])
	case "repl":
		err = repl.RunREPL(ctx, config)
	case "web":
		err = web.Run(ctx, config)
	default:
		err = fmt.Errorf("unknown command: %s", os.Args[1])
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
