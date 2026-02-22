package task

import (
	"context"
	"fmt"
	"time"

	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// RunTask orchestrates the complete task execution flow:
// 1. Setup the task
// 2. Show task intro screen
// 3. Run the interface
// 4. Start validation loop in background
// 5. Show questionnaire when done
func RunTask(ctx context.Context, t Tasker, i Interface, ifaceType InterfaceType) error {
	doneChan := make(chan bool, 1)
	quitChan := make(chan bool, 1)
	feedbackChan := make(chan ValidationFeedback, 10)

	if validated, ok := i.(ValidatedInterface); ok {
		validated.SetChannels(feedbackChan, quitChan)
	}

	// Setup task
	if err := t.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup task: %w", err)
	}

	// Show task intro
	detailsScreen := taskui.NewTaskModel(t.Details())
	model, err := tea.NewProgram(detailsScreen, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(interface{ GetUserQuit() bool }); ok && m.GetUserQuit() {
		return ErrUserQuit
	}

	// Start validation loop
	go startValidationLoop(ctx, t, feedbackChan, doneChan, quitChan)

	// Run interface
	interfaceDone := make(chan error, 1)
	go func() {
		interfaceDone <- i.Run(ctx, t.Config())
	}()

	select {
	case <-doneChan:
		close(quitChan)
		if err := <-interfaceDone; err != nil {
			fmt.Printf("warning: interface error after task completion: %v\n", err)
		}
		fmt.Println("Task completed successfully!")

	case err := <-interfaceDone:
		close(quitChan)
		if err != nil {
			return fmt.Errorf("failed to start task interface: %w", err)
		}
		fmt.Println("Task incomplete - you exited early")
	}

	// Show questionnaire
	questions := t.Questions(ifaceType)
	questionare := taskui.NewQuestionnaireModel(questions)
	model, err = tea.NewProgram(questionare, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(interface{ GetUserQuit() bool }); ok && m.GetUserQuit() {
		return ErrUserQuit
	}

	return nil
}

func startValidationLoop(ctx context.Context, t Tasker, feedbackChan chan ValidationFeedback, doneChan chan bool, quitChan chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ok, err := t.Validate(ctx)
			feedback := ValidationFeedback{
				Success: ok,
			}
			if ok {
				feedback.Message = "Task completed successfully!"
				feedbackChan <- feedback
				doneChan <- true
				return
			}
			if err != nil {
				feedback.Message = err.Error()
			} else {
				feedback.Message = "Task not yet complete"
			}
			feedbackChan <- feedback
		case <-quitChan:
			return
		case <-ctx.Done():
			return
		}
	}
}
