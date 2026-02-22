// Package cli provides the command-line interface for the PM System.
package cli

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/commands/issues"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

// Config is an alias for service.Config, used to configure the CLI.
type Config = service.Config

type CLI struct {
	RootCmd *cobra.Command
}

func NewCli(rootCmd *cobra.Command) *CLI {
	return &CLI{
		RootCmd: rootCmd,
	}
}

// Run initializes the services and executes the CLI commands.
func (c *CLI) Run(ctx context.Context, config Config) error {
	app, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	issuesCmd.SetApp(app)

	if err := fang.Execute(ctx, c.RootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme)); err != nil {
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

	issuesCmd.SetApp(app)

	if err := issuesCmd.ExecuteArgs(args); err != nil {
		return err
	}

	return nil
}
