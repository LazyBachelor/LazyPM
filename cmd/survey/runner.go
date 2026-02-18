package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
	"github.com/LazyBachelor/LazyPM/pkg/task"
)

func runTask(ctx context.Context, t task.Tasker, i task.Interface) error {
	return task.RunTask(ctx, t, i, tasks.InterfaceToType(i))
}

func taskLoop(ctx context.Context, surveyTasks []task.Tasker, interfaces map[string]task.Interface) error {
	var ifaceNames []string
	for name := range interfaces {
		ifaceNames = append(ifaceNames, name)
	}

	rand.Shuffle(len(ifaceNames), func(i, j int) {
		ifaceNames[i], ifaceNames[j] = ifaceNames[j], ifaceNames[i]
	})

	for i, t := range surveyTasks {
		idx := i % len(ifaceNames)
		selected := interfaces[ifaceNames[idx]]

		if err := runTask(ctx, t, selected); err != nil {
			return err
		}
	}
	return nil
}

func returnIfUserQuit(err error, msg string) error {
	if errors.Is(err, task.ErrUserQuit) {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
