package task

import (
	"context"
	"fmt"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

type Task struct {
	Config        TaskConfig
	interfaceType Interface
	aboutScreen   tea.Model
	questionnaire tea.Model

	validateFunc ValidateFunc
	dbStateFunc  DbStateFunc

	svc *service.Services

	feedbackChan chan ValidationFeedback
	doneChan     chan bool
	quitChan     chan bool
}

func NewTask(svc *service.Services, aboutScreen tea.Model, questionnaire tea.Model) *Task {
	return &Task{
		aboutScreen:   aboutScreen,
		questionnaire: questionnaire,
		svc:           svc,
	}
}

func (t *Task) IntroduceTask() error {
	if t.aboutScreen == nil {
		return fmt.Errorf("aboutScreen is not set")
	}
	model, err := tea.NewProgram(t.aboutScreen, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(interface{ GetUserQuit() bool }); ok && m.GetUserQuit() {
		return ErrUserQuit
	}
	return nil
}

func (t *Task) StartInterface(ctx context.Context, cfg TaskConfig) error {
	if t.interfaceType == nil {
		return fmt.Errorf("interfaceType is not set")
	}

	return t.interfaceType.Run(ctx, cfg)
}

func (t *Task) Initialize(ctx context.Context) error {
	if t.dbStateFunc == nil {
		return fmt.Errorf("dbStateFunc is not set")
	}
	return t.dbStateFunc(ctx, t.svc)
}

func (t *Task) Validate(ctx context.Context) (bool, error) {
	if t.validateFunc == nil {
		return false, fmt.Errorf("validateFunc is not set")
	}
	return t.validateFunc(ctx, t.svc)
}

func (t *Task) StartQuestionnaire() error {
	if t.questionnaire == nil {
		return fmt.Errorf("questionnaire is not set")
	}
	model, err := tea.NewProgram(t.questionnaire, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(interface{ GetUserQuit() bool }); ok && m.GetUserQuit() {
		return ErrUserQuit
	}
	return nil
}

func (t *Task) SetConfigFunc(fn ConfigFunc) {
	t.Config = fn()
}

func (t *Task) SetInterface(interfaceType Interface) {
	t.interfaceType = interfaceType
}

func (t *Task) SetDbStateFunc(fn DbStateFunc) {
	t.dbStateFunc = fn
}

func (t *Task) SetValidateFunc(fn ValidateFunc) {
	t.validateFunc = fn
}

func (t *Task) SetChannels(feedbackChan chan ValidationFeedback, doneChan chan bool, quitChan chan bool) {
	t.feedbackChan = feedbackChan
	t.doneChan = doneChan
	t.quitChan = quitChan
}

func (t *Task) StartValidationLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ok, err := t.Validate(ctx)
			feedback := ValidationFeedback{
				Timestamp: time.Now(),
			}
			if ok {
				feedback.Success = true
				feedback.Message = "Task completed successfully!"
				t.feedbackChan <- feedback
				t.doneChan <- true
				return
			} else {
				feedback.Success = false
				if err != nil {
					feedback.Message = err.Error()
				} else {
					feedback.Message = "Task not yet complete"
				}
				t.feedbackChan <- feedback
			}
		case <-t.quitChan:
			return
		case <-ctx.Done():
			return
		}
	}
}
