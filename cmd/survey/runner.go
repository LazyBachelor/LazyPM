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

		doneChan := make(chan bool, 1)
		quitChan := make(chan bool, 1)
		feedbackChan := make(chan task.ValidationFeedback, 10)

		if validated, ok := interfaces[interfaceIndex].(task.ValidatedInterface); ok {
			validated.SetChannels(feedbackChan, quitChan)
		}

		if err := t.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize task: %w", err)
		}

		if err := t.IntroduceTask(); err != nil {
			return returnIfUserQuit(err, "failed to display task introduction screen")
		}

		t.SetChannels(feedbackChan, doneChan, quitChan)

		go t.StartValidationLoop(ctx)

		interfaceDone := make(chan error, 1)
		go func() {
			interfaceDone <- t.StartInterface(ctx, t.Config)
		}()

		select {
		case <-doneChan:
			close(quitChan)
			<-interfaceDone
			fmt.Println("Task completed successfully!")

		case err := <-interfaceDone:
			close(quitChan)
			if err != nil {
				return returnIfUserQuit(err, "failed to start task interface")
			}
			fmt.Println("Task incomplete - you exited early")
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
