package task

import (
	"context"
	"fmt"

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
}

func NewTask(aboutScreen tea.Model, questionnaire tea.Model) *Task {
	return &Task{
		aboutScreen:   aboutScreen,
		questionnaire: questionnaire,
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

func (t *Task) Initialize(ctx context.Context, svc *service.Services) error {
	if t.dbStateFunc == nil {
		return fmt.Errorf("dbStateFunc is not set")
	}
	return t.dbStateFunc(ctx, svc)
}

func (t *Task) Validate(ctx context.Context, svc *service.Services) (bool, error) {
	if t.validateFunc == nil {
		return false, fmt.Errorf("validateFunc is not set")
	}
	return t.validateFunc(ctx, svc)
}
