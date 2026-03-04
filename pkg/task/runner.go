package task

import (
	"context"
	"fmt"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

type App = models.App
type Config = models.Config

type Tasker = models.Tasker

type Interface = models.Interface
type InterfaceType = models.InterfaceType

type ValidatedInterface = models.ValidatedInterface
type ValidationFeedback = models.ValidationFeedback

var ErrUserQuit = models.ErrUserQuit

// RunTask orchestrates the complete task execution flow:
// 1. Setup the task
// 2. Show task intro screen
// 3. Run the interface
// 4. Start validation loop in background
// 5. Show questionnaire when done
func RunTask(ctx context.Context, t Tasker, i Interface, iType InterfaceType) error {
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
	detailsScreen := NewTaskModel(t.Details())
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
	questions := t.Questions(iType)
	questionare := NewQuestionnaireModel(questions)
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
			feedback := t.Validate(ctx)
			if feedback.Success {
				if feedback.Message == "" {
					feedback.Message = "Task completed successfully! Going back to the survey menu..."
				}
				feedbackChan <- feedback
				time.Sleep(4 * time.Second)
				doneChan <- true
				return
			}
			feedback.Message = "Task not completed!"
			feedbackChan <- feedback
		case <-quitChan:
			return
		case <-ctx.Done():
			return
		}
	}
}
