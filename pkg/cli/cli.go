package cli

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli/commands"
	"context"
)

type CLIConfig = service.Config

func Run(ctx context.Context, config CLIConfig) error {
	svc, cleanup, err := service.NewServices(ctx, config)

	if err != nil {
		return err
	}

	defer cleanup()

	if err := commands.Execute(svc); err != nil {
		return err
	}

	return nil
}
