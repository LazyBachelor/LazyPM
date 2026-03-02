package task

import (
	"context"
	"fmt"
	"log/slog"
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

type QuestionnaireKeysProvider interface {
	QuestionnaireKeys(InterfaceType) []string
}

var ErrUserQuit = models.ErrUserQuit

// RunTask orchestrates the complete task execution flow:
// 1. Setup the task
// 2. Show task intro screen
// 3. Run the interface
// 4. Start validation loop in background
// 5. Show questionnaire when done
func RunTask(ctx context.Context, app *App, t Tasker, i Interface, iType InterfaceType) (runErr error) {
	config := t.Config()
	details := t.Details()

	var logger *slog.Logger
	if app != nil {
		logger = app.Logger
	}

	collector := newTaskRunCollector(details.Title, iType, logger)
	collector.log("info", "task run started")

	config = config.WithActionLogger(func(action string) {
		collector.recordUserAction(action)
	})

	defer func() {
		if runErr != nil {
			collector.setError(runErr)
			collector.log("error", runErr.Error())
		}

		run := collector.finalize()

		if err := appendTaskMetrics(config.StatisticsStoragePath, details.Title, run, logger); err != nil {
			if runErr == nil {
				runErr = fmt.Errorf("failed to persist task metrics: %w", err)
				return
			}
			if logger != nil {
				logger.Warn("failed to persist task metrics", "error", err, "task", details.Title)
			}
		}

		if app != nil && app.Stats != nil {
			if err := app.Stats.RecordTaskRun(ctx, run); err != nil {
				if runErr == nil {
					runErr = fmt.Errorf("failed to update global statistics: %w", err)
					return
				}
				if logger != nil {
					logger.Warn("failed to update global statistics", "error", err, "task", details.Title)
				}
			}
		}
	}()

	doneChan := make(chan bool, 1)
	quitChan := make(chan bool, 1)
	feedbackChan := make(chan ValidationFeedback, 10)

	if validated, ok := i.(ValidatedInterface); ok {
		validated.SetChannels(feedbackChan, quitChan)
	}

	// Setup task
	collector.log("info", "setting up task")
	if err := t.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup task: %w", err)
	}
	collector.log("info", "task setup completed")

	// Show task intro
	collector.log("info", "showing task intro")
	detailsScreen := NewTaskModel(details)
	model, err := tea.NewProgram(detailsScreen, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(interface{ GetUserQuit() bool }); ok && m.GetUserQuit() {
		collector.log("info", "user quit during task intro")
		return ErrUserQuit
	}
	collector.log("info", "task intro completed")

	// Start validation loop
	collector.log("info", "starting validation loop")
	go startValidationLoop(ctx, t, feedbackChan, doneChan, quitChan, collector.recordValidation)

	// Run interface
	collector.log("info", "starting task interface")
	interfaceDone := make(chan error, 1)
	go func() {
		interfaceDone <- i.Run(ctx, config)
	}()

	select {
	case <-doneChan:
		close(quitChan)
		collector.setCompleted(true)
		collector.log("info", "task validation completed")
		if err := <-interfaceDone; err != nil {
			if logger != nil {
				logger.Warn("interface error after task completion", "error", err, "task", details.Title)
			}
			collector.log("warn", fmt.Sprintf("interface error after task completion: %v", err))
		}
		fmt.Println("Task completed successfully!")

	case err := <-interfaceDone:
		close(quitChan)
		if err != nil {
			collector.log("error", fmt.Sprintf("task interface failed: %v", err))
			return fmt.Errorf("failed to start task interface: %w", err)
		}
		collector.log("info", "task interface exited before completion")
		fmt.Println("Task incomplete - you exited early")
	}

	// Show questionnaire
	collector.log("info", "showing post-task questionnaire")
	questions := t.Questions(iType)
	questionnaireKeys := []string{}
	if provider, ok := t.(QuestionnaireKeysProvider); ok {
		questionnaireKeys = provider.QuestionnaireKeys(iType)
	}
	questionare := NewQuestionnaireModel(questions, questionnaireKeys)
	model, err = tea.NewProgram(questionare, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}

	questionnaireCompleted := false
	questionnaireAnswers := map[string]any(nil)
	if m, ok := model.(interface {
		GetCompleted() bool
		GetAnswers() map[string]any
	}); ok {
		questionnaireCompleted = m.GetCompleted()
		questionnaireAnswers = m.GetAnswers()
	}

	if m, ok := model.(interface{ GetUserQuit() bool }); ok && m.GetUserQuit() {
		collector.recordQuestionnaire(questionnaireCompleted, true, questionnaireAnswers)
		collector.log("info", "user quit during questionnaire")
		return ErrUserQuit
	}
	collector.recordQuestionnaire(questionnaireCompleted, false, questionnaireAnswers)
	collector.log("info", "task run finished")

	return nil
}

func startValidationLoop(ctx context.Context, t Tasker, feedbackChan chan ValidationFeedback, doneChan chan bool, quitChan chan bool, onFeedback func(ValidationFeedback)) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			feedback := t.Validate(ctx)
			if onFeedback != nil {
				onFeedback(feedback)
			}
			if feedback.Success {
				feedback.Message = "Task completed successfully!"
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
