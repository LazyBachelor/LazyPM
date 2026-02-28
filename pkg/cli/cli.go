// Package cli provides the command-line interface for the PM System.
package cli

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/app"
	issues "github.com/LazyBachelor/LazyPM/internal/commands/issues"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

// Config is an alias for app.Config, used to configure the CLI.
type Config = models.Config

type CLI struct {
	RootCmd *cobra.Command
}

func New(rootCmd *cobra.Command) *CLI {
	return &CLI{
		RootCmd: rootCmd,
	}
}

// Run initializes the services and executes the CLI commands.
func (c *CLI) Run(ctx context.Context, config Config) error {
	app, cleanup, err := app.New(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	issues.SetApp(app)

	if err := fang.Execute(ctx, c.RootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme)); err != nil {
		return err
	}

	return nil
}

// RunWithArgs initializes the services and executes the CLI commands with the provided arguments.
func (c *CLI) RunWithArgs(ctx context.Context, config Config, args []string) error {
	app, cleanup, err := app.New(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	issues.SetApp(app)

	if err := issues.ExecuteArgs(args); err != nil {
		return err
	}

	return nil
}
