package cli

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli/commands"
)

type CLIConfig = service.Config

func Run(ctx context.Context, config CLIConfig) error {
	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	commands.SetServices(svc)

	if err := commands.Execute(); err != nil {
		return err
	}

	return nil
}

func RunWithArgs(ctx context.Context, config CLIConfig, args []string) error {
	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	commands.SetServices(svc)

	if err := commands.ExecuteArgs(args); err != nil {
		return err
	}

	return nil
}
