package task

import (
	"context"
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/service"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
)

type Config = service.Config
type InterfaceType string

type Interface interface {
	Run(context.Context, Config) error
}

type Tasker interface {
	Config() Config
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

var ErrUserQuit = fmt.Errorf("user quit")
