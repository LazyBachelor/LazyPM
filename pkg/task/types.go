package task

import (
	"context"
	"errors"

	"github.com/LazyBachelor/LazyPM/internal/service"
)

var ErrUserQuit = errors.New("user quit")

type TaskConfig = service.Config

type Interface interface {
	Run(context.Context, TaskConfig) error
}

type ConfigFunc func() TaskConfig
type ValidateFunc func(context.Context, *service.Services) (ok bool, err error)
type DbStateFunc func(context.Context, *service.Services) error
