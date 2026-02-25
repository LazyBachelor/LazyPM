package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"

	_ "github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
)

func initializeServices(ctx context.Context) (*service.App, func(), error) {
	return service.NewServices(ctx, tasks.BaseConfig())
}

func initInterfaces() map[string]task.Interface {
	interfaces := make(map[string]task.Interface)
	for _, name := range task.ListInterfaces() {
		i, err := task.GetInterface(name)
		if err != nil {
			continue
		}
		interfaces[name] = i
	}
	return interfaces
}

func initTasks(app *service.App) map[string]task.Tasker {
	taskMap := make(map[string]task.Tasker)
	for _, name := range task.ListTasks() {
		t, err := task.GetTask(name, app)
		if err != nil {
			continue
		}
		taskMap[name] = t
	}

	return taskMap
}
