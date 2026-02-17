package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
)

func taskLoop(ctx context.Context, svc *service.Services, surveyTasks []*task.Task, interfaces []task.Interface) error {
	interfaceIndex := rand.Int() % len(interfaces)

	for _, t := range surveyTasks {

		t.SetInterface(interfaces[interfaceIndex])

		if err := t.Initialize(ctx, svc); err != nil {
			return fmt.Errorf("failed to initialize task: %w", err)
		}

		if err := t.IntroduceTask(); err != nil {
			if errors.Is(err, task.ErrUserQuit) {
				return task.ErrUserQuit
			}
			return fmt.Errorf("failed to display task introduction screen: %w", err)
		}

		if err := t.StartInterface(ctx, t.Config); err != nil {
			return fmt.Errorf("failed to start task interface: %w", err)
		}

		ok, err := t.Validate(ctx, svc)
		if err != nil {
			return fmt.Errorf("validation error: %w", err)
		}
		if !ok {
			return fmt.Errorf("task validation failed: task did not meet requirements")
		}

		if err := t.StartQuestionnaire(); err != nil {
			if errors.Is(err, task.ErrUserQuit) {
				return task.ErrUserQuit
			}
			return fmt.Errorf("failed to start questionnaire: %w", err)
		}

		interfaceIndex++
		if interfaceIndex >= len(interfaces) {
			interfaceIndex = 0
		}
	}
	return nil
}
