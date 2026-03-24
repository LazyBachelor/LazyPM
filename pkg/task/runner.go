package task

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"golang.org/x/term"
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

type TaskRunner struct {
	app    *App
	logger *slog.Logger
}

func NewTaskRunner(app *App) *TaskRunner {
	var logger *slog.Logger
	if app != nil {
		logger = app.Logger
	}
	return &TaskRunner{
		app:    app,
		logger: logger,
	}
}

func (r *TaskRunner) Run(ctx context.Context, t Tasker, i Interface, iType InterfaceType) (runErr error) {

	config := t.Config()
	details := t.Details(iType)

	lifecycle := NewRunLifecycle(r.app, config, details, iType, r.logger)

	defer func() {
		runErr = lifecycle.Finish(ctx, runErr)
	}()

	collector := lifecycle.collector
	config = lifecycle.config

	collector.log("info", "task run started")

	// Setup
	if err := t.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup task: %w", err)
	}

	// Intro screen
	if err := runIntro(details); err != nil {
		return err
	}

	// Validation
	feedbackChan := make(chan ValidationFeedback, 10)
	quitChan := make(chan bool, 1)
	submitChan := make(chan models.ValidationTrigger, 1)

	if validated, ok := i.(ValidatedInterface); ok {
		validated.SetChannels(feedbackChan, quitChan)
		validated.SetSubmitChan(submitChan)
	}

	if r.app != nil {
		r.app.SubmitChan = submitChan
	}

	engine := &ValidationEngine{task: t}
	doneChan, stopChan := engine.Start(ctx, submitChan, func(feedback ValidationFeedback, source models.ValidationTriggerSource) {
		collector.recordValidation(feedback, source)

		if feedback.Success {
			collector.setCompleted(true)
			feedback.Message = "Task completed successfully!"
		} else if feedback.Message == "" {
			feedback.Message = "Task not completed!"
		}

		if r.app != nil {
			r.app.CurrentFeedback = &feedback
		}

		select {
		case feedbackChan <- feedback:
		default:
		}
	})

	// Run interface
	interfaceErr := make(chan error, 1)
	go func() {
		interfaceErr <- i.Run(ctx, config)
	}()

	select {
	case <-doneChan:
		close(stopChan)
		close(quitChan)
		collector.setCompleted(true)
		if err := <-interfaceErr; err != nil {
			return fmt.Errorf("task interface failed during shutdown: %w", err)
		}

	case err := <-interfaceErr:
		close(stopChan)
		close(quitChan)
		if err != nil {
			return fmt.Errorf("task interface failed: %w", err)
		}
	}

	if oldState, err := term.GetState(int(os.Stdin.Fd())); err == nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}

	fmt.Print("\033[0m\033[?25h")

	// Questionnaire
	if err := runQuestionnaire(t, iType, collector); err != nil {
		return err
	}

	collector.log("info", "task run finished")
	return nil
}

func (r *RunLifecycle) Finish(ctx context.Context, runErr error) error {

	if runErr != nil {
		r.collector.setError(runErr)
		r.collector.log("error", runErr.Error())
	}

	run := r.collector.finalize()

	if r.metricsStore != nil {
		if err := r.metricsStore.Append(ctx, r.details.Title, run); err != nil && runErr == nil {
			return fmt.Errorf("persist metrics: %w", err)
		}
	}

	if r.app != nil && r.app.Stats != nil {
		if err := r.app.Stats.RecordTaskRun(ctx, run); err != nil && runErr == nil {
			return fmt.Errorf("update global stats: %w", err)
		}
	}

	return runErr
}

func runIntro(details models.TaskDetails) error {
	model, err := tea.NewProgram(NewTaskModel(details)).Run()
	if err != nil {
		return err
	}

	if m, ok := model.(interface{ GetUserQuit() bool }); ok {
		if m.GetUserQuit() {
			return ErrUserQuit
		}
	}

	return nil
}

func runQuestionnaire(t Tasker, iType InterfaceType, collector *taskRunCollector) error {
	questions := t.Questions(iType)

	keys := []string{}
	if provider, ok := t.(QuestionnaireKeysProvider); ok {
		keys = provider.QuestionnaireKeys(iType)
	}

	model, err := tea.NewProgram(NewQuestionnaireModel(questions, keys)).Run()
	if err != nil {
		return err
	}

	var completed bool
	var answers map[string]any

	if m, ok := model.(interface {
		GetCompleted() bool
		GetAnswers() map[string]any
	}); ok {
		completed = m.GetCompleted()
		answers = m.GetAnswers()
	}

	userQuit := false
	if m, ok := model.(interface{ GetUserQuit() bool }); ok {
		userQuit = m.GetUserQuit()
	}

	collector.recordQuestionnaire(completed, userQuit, answers)

	if userQuit {
		return ErrUserQuit
	}

	return nil
}
