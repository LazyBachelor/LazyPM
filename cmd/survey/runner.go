package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
	surveyCmd "github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/spf13/cobra"
)

func runStartCmd(cmd *cobra.Command, args []string) error {
	interfaces := initInterfaces()

	svc, cleanup, err := initializeServices(cmd.Context())
	if err != nil {
		return returnIfUserQuit(err, "failed to initialize services")
	}
	defer cleanup()

	surveyTasks := initTasks(svc)

	if cmd.Flags().Changed("interface") {
		if _, ok := interfaces[surveyCmd.InterfaceType]; !ok {
			return fmt.Errorf("invalid interface, valid are %v", task.ListInterfaces())
		}
		interfaces = map[string]task.Interface{
			surveyCmd.InterfaceType: interfaces[surveyCmd.InterfaceType],
		}
	}

	if cmd.Flags().Changed("stage") {
		if surveyCmd.Task < 1 || surveyCmd.Task > len(surveyTasks) {
			return fmt.Errorf("invalid stage")
		}
		if err := task.RunTask(cmd.Context(), surveyTasks[surveyCmd.Task-1],
			interfaces[surveyCmd.InterfaceType], tasks.InterfaceToType(interfaces[surveyCmd.InterfaceType])); err != nil {
			return err
		}
		return nil
	}

	if err := newIntroModel().Run(); err != nil {
		return returnIfUserQuit(err, "failed to run intro")
	}

	if err := taskLoop(cmd.Context(), surveyTasks, interfaces); err != nil {
		return returnIfUserQuit(err, "task loop failed")
	}
	return nil
}

func taskLoop(ctx context.Context, surveyTasks []task.Tasker, interfaces map[string]task.Interface) error {
	var iNames []string
	for name := range interfaces {
		iNames = append(iNames, name)
	}

	rand.Shuffle(len(iNames), func(i, j int) {
		iNames[i], iNames[j] = iNames[j], iNames[i]
	})

	for i, t := range surveyTasks {
		idx := i % len(iNames)
		selected := interfaces[iNames[idx]]

		if err := task.RunTask(ctx, t, selected, tasks.InterfaceToType(selected)); err != nil {
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
