package task

import (
	"context"
	"errors"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/service"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
)

var ErrUserQuit = errors.New("user quit")

type TaskConfig = service.Config

type Interface interface {
	Run(context.Context, TaskConfig) error
}

type InterfaceType string

const (
	InterfaceTUI InterfaceType = "tui"
	InterfaceCLI InterfaceType = "repl"
	InterfaceWeb InterfaceType = "web"
)

type ConfigFunc func() TaskConfig
type ValidateFunc func(context.Context) (ok bool, err error)
type DbStateFunc func(context.Context) error
type QuestionsFunc func(InterfaceType) taskui.Questions

type ValidationFeedback struct {
	Success   bool
	Message   string
	Timestamp time.Time
}

type ValidationObserver interface {
	OnValidationUpdate(feedback ValidationFeedback)
	OnTaskComplete()
}

type ValidatedInterface interface {
	Interface
	SetChannels(feedbackChan chan ValidationFeedback, quitChan chan bool)
}
