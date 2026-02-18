package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
)

func runTask(ctx context.Context, t task.Tasker, i task.Interface) error {
	return task.RunTask(ctx, t, i, InterfaceToType(i))
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

func InterfaceToType(it task.Interface) task.InterfaceType {
	switch it.(type) {
	case *repl.REPL:
		return task.InterfaceCLI
	case *tui.Tui:
		return task.InterfaceTUI
	case *web.Web:
		return task.InterfaceWeb
	default:
		return task.InterfaceType("unknown")
	}
}

func returnIfUserQuit(err error, msg string) error {
	if errors.Is(err, task.ErrUserQuit) {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
