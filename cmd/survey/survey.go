package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
)

func main() {
	ctx := context.Background()

	svc, close, err := initializeServices(ctx)
	if err != nil {
		fatal("Failed to initialize services: %v\n", err)
	}
	defer close()

	tui := tui.NewTui()
	web := web.NewWeb()
	repl := repl.NewRepl()

	createIssueTask := tasks.NewCreateIssueTask(web)

	task := createIssueTask

	task.SetInterface(tui)
	task.SetInterface(repl)

	task.MigrateToTask(ctx, svc, task)

	if err := task.IntroduceTask(); err != nil {
		fatal("Failed to introduce task: %v\n", err)
	}

	if err := task.StartInterface(ctx, svc.Config); err != nil {
		fatal("Failed to start interface: %v\n", err)
	}

	if err := task.StartQuestionnaire(); err != nil {
		fatal("Failed to start questionnaire: %v\n", err)
	}
}

func initializeServices(ctx context.Context) (*service.Services, func(), error) {
	config := service.Config{
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/stats.json",
		WebAddress:            "localhost:8080",
	}
	return service.NewServices(ctx, config)
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
