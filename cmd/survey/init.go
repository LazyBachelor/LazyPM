package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
)

func initializeServices(ctx context.Context) (*service.Services, func(), error) {
	config := service.Config{
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/stats.json",
		WebAddress:            "localhost:8080",
	}
	return service.NewServices(ctx, config)
}

func initTasks(svc *service.Services) []*task.Task {
	return []*task.Task{
		tasks.NewCreateIssueTask(svc).Init(),
	}
}

func initInterfaces() []task.Interface {
	return []task.Interface{repl.NewRepl(), tui.NewTui(), web.NewWeb()}
}
