package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
)

type TaskConfig = service.Config

type Interface interface {
	Run(context.Context, TaskConfig) error
}

type ConfigFunc func() TaskConfig
type ValidateFunc func(context.Context, *service.Services) (ok bool, err error)
type DbStateFunc func(context.Context, *service.Services) error
