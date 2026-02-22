package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"

	_ "github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
)

func initializeServices(ctx context.Context) (*service.App, func(), error) {
	return service.NewServices(ctx, tasks.BaseConfig())
}

func initInterfaces() map[string]task.Interface {
	return map[string]task.Interface{
		"repl": repl.NewRepl(),
		"tui":  tui.NewTui(),
		"web":  web.NewWeb(),
	}
}

func initTasks(app *service.App) []task.Tasker {
	var taskList []task.Tasker

	for _, name := range task.List() {
		t, err := task.Get(name, app)
		if err != nil {
			continue
		}
		taskList = append(taskList, t)
	}

	return taskList
}
