// Package cli provides the command-line interface for the PM System.
package cli

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli/commands"
)

// Config is an alias for service.Config, used to configure the CLI.
type Config = service.Config

type CLI struct{}

func NewCli() *CLI {
	return &CLI{}
}

// Run initializes the services and executes the CLI commands.
func (c *CLI) Run(ctx context.Context, config Config) error {
	app, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	commands.SetApp(app)

	if err := commands.Execute(); err != nil {
		return err
	}

	return nil
}

// RunWithArgs initializes the services and executes the CLI commands with the provided arguments.
func (c *CLI) RunWithArgs(ctx context.Context, config Config, args []string) error {
	app, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	commands.SetApp(app)

	if err := commands.ExecuteArgs(args); err != nil {
		return err
	}

	return nil
}
