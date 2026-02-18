package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"

	_ "github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
)

func initializeServices(ctx context.Context) (*service.Services, func(), error) {
	config := service.Config{
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/stats.json",
		WebAddress:            ":8080",
	}
	return service.NewServices(ctx, config)
}

func initInterfaces() map[string]task.Interface {
	return map[string]task.Interface{
		"repl": repl.NewRepl(),
		"tui":  tui.NewTui(),
		"web":  web.NewWeb(),
	}
}

func initTasks(svc *service.Services) []task.Tasker {
	var taskers []task.Tasker

	for _, name := range task.List() {
		t, err := task.Get(name, svc)
		if err != nil {
			continue
		}
		taskers = append(taskers, t)
	}

	return taskers
}
