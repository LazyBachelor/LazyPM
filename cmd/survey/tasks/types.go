package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

type Interface interface {
	Run(context.Context, service.Config) error
}

type ValidateFunc func(context.Context, *service.Services) (ok bool, err error)
type DbStateFunc func(context.Context, *service.Services) error

type Task struct {
	interfaceType Interface
	aboutScreen   tea.Model
	questionnaire tea.Model

	validateFunc ValidateFunc
	dbStateFunc  DbStateFunc
}

type TaskList struct {
	Todo []*Task
	Done []*Task
}
