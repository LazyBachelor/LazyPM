package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
)

func main() {
	ctx := context.Background()

	svc, close, err := initializeServices(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v\n", err)
	}
	defer close()

	surveyTasks := initTasks()
	interfaces := initInterfaces()

	if err := taskLoop(ctx, svc, surveyTasks, interfaces); err != nil {
		log.Fatalf("Task loop failed: %v\n", err)
	}
}

func taskLoop(ctx context.Context, svc *service.Services, surveyTasks []*task.Task, interfaces []task.Interface) error {
	interfaceIndex := rand.Int() % len(interfaces)

	for _, task := range surveyTasks {

		task.SetInterface(interfaces[interfaceIndex])

		if err := task.Initialize(ctx, svc); err != nil {
			return fmt.Errorf("failed to initialize task: %w", err)
		}

		if err := task.IntroduceTask(); err != nil {
			return fmt.Errorf("failed to display task introduction screen: %w", err)
		}

		if err := task.StartInterface(ctx, task.Config); err != nil {
			return fmt.Errorf("failed to start task interface: %w", err)
		}

		ok, err := task.Validate(ctx, svc)
		if err != nil {
			return fmt.Errorf("validation error: %w", err)
		}
		if !ok {
			return fmt.Errorf("task validation failed: task did not meet requirements")
		}

		if err := task.StartQuestionnaire(); err != nil {
			return fmt.Errorf("failed to start questionnaire: %w", err)
		}

		interfaceIndex++
		if interfaceIndex >= len(interfaces) {
			interfaceIndex = 0
		}
	}
	return nil
}
