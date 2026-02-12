// Package cli provides the command-line interface for the PM System.
package cli

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli/commands"
)

// CLIConfig is an alias for service.Config, used to configure the CLI.
type CLIConfig = service.Config

type CLI struct{}

func NewCli() *CLI {
	return &CLI{}
}

// Run initializes the services and executes the CLI commands.
func (c *CLI) Run(ctx context.Context, config CLIConfig) error {
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

// RunWithArgs initializes the services and executes the CLI commands with the provided arguments.
func (c *CLI) RunWithArgs(ctx context.Context, config CLIConfig, args []string) error {
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
