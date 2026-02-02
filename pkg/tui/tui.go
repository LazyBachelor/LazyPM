package tui

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"context"
)

type TUIConfig = service.Config

func Run(ctx context.Context, config TUIConfig) error {
	_, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}
	defer cleanup()

	return nil
}
