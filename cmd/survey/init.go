package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
	_ "github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
)

func init() {
	task.Register("create_issue", func(svc *service.Services) task.Tasker {
		return tasks.NewCreateIssueTask(svc)
	})
	task.Register("coding_task", func(svc *service.Services) task.Tasker {
		return tasks.NewCodingTask(svc)
	})
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

func initializeServices(ctx context.Context) (*service.Services, func(), error) {
	return service.NewServices(ctx, tasks.BaseConfig())
}

func initInterfaces() map[string]task.Interface {
	return map[string]task.Interface{
		"repl": repl.NewRepl(),
		"tui":  tui.NewTui(),
		"web":  web.NewWeb(),
	}
}
