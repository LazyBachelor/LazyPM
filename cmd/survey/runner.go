package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/LazyBachelor/LazyPM/pkg/task"
)

func taskLoop(ctx context.Context, surveyTasks []*task.Task, interfaces []task.Interface) error {
	interfaceIndex := rand.Int() % len(interfaces)

	for _, t := range surveyTasks {

		t.SetInterface(interfaces[interfaceIndex])

		if err := t.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize task: %w", err)
		}

		if err := t.IntroduceTask(); err != nil {
			return returnIfUserQuit(err, "failed to display task introduction screen")
		}

		if err := t.StartInterface(ctx, t.Config); err != nil {
			return returnIfUserQuit(err, "failed to start task interface")
		}

		ok, err := t.Validate(ctx)
		if err != nil {
			return returnIfUserQuit(err, "validation error")
		}
		if !ok {
			return fmt.Errorf("task validation failed: task did not meet requirements")
		}

		if err := t.StartQuestionnaire(); err != nil {
			return returnIfUserQuit(err, "failed to start questionnaire")
		}

		interfaceIndex++
		if interfaceIndex >= len(interfaces) {
			interfaceIndex = 0
		}
	}
	return nil
}

func returnIfUserQuit(err error, msg string) error {
	if errors.Is(err, ErrUserQuit) {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
