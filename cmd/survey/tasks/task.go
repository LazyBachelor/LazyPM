// task.go
package tasks

import (
	"context"
	"errors"

	"github.com/LazyBachelor/LazyPM/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

func NewTask(interfaceType Interface, aboutScreen tea.Model, questionnaire tea.Model) *Task {
	return &Task{
		interfaceType: interfaceType,
		aboutScreen:   aboutScreen,
		questionnaire: questionnaire,
	}
}

func (t *Task) SetValidateFunc(fn ValidateFunc) {
	t.validateFunc = fn
}

func (t *Task) SetDbStateFunc(fn DbStateFunc) {
	t.dbStateFunc = fn
}

func (t *Task) SetInterface(interfaceType Interface) {
	t.interfaceType = interfaceType
}

func (t *Task) Validate(ctx context.Context, svc *service.Services) (bool, error) {
	if t.validateFunc == nil {
		return false, errors.New("validateFunc is not set")
	}
	return t.validateFunc(ctx, svc)
}

func (t *Task) MigrateToTask(ctx context.Context, svc *service.Services, task *Task) error {
	t = task
	if t.dbStateFunc == nil {
		return errors.New("dbStateFunc is not set")
	}
	return t.dbStateFunc(ctx, svc)
}
