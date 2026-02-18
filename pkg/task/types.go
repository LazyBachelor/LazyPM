package task

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
)

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

type Tasker interface {
	Config() TaskConfig
	Details() taskui.TaskDetails
	Questions(InterfaceType) taskui.Questions
	Setup(context.Context) error
	Validate(context.Context) (bool, error)
}

type ValidationFeedback struct {
	Success bool
	Message string
}

type ValidatedInterface interface {
	Interface
	SetChannels(feedbackChan chan ValidationFeedback, quitChan chan bool)
}
