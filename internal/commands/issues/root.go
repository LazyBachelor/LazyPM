package issues

import (
	"bytes"
	"context"

	"charm.land/fang/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"

	"github.com/spf13/cobra"
)

type contextKey string

const appKey contextKey = "app"

type App = models.App

var app *App

// Flags struct to hold command-line flag values for issues.
type Flags struct {
	interactive bool
	limit       int

	title       string
	description string
	status      string
	issueType   string
	priority    int
	assignee    string
}

// RootCmd is the base command for the CLI application.
var RootCmd = &cobra.Command{
	Use:   "pm",
	Short: "Project Management User Interface Comparison CLI",
	Long:  `Project Management User Interface Comparison CLI is a tool designed to evaluate and compare different project management interfaces through a series of tasks and surveys.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Inject app into context for all commands
		if app != nil {
			cmd.SetContext(context.WithValue(cmd.Context(), appKey, app))
		}
	},
}

// SetApp sets the app variable for use in command execution.
// Must be called before executing any commands to ensure services are available.
func SetApp(application *App) {
	app = application
	RootCmd.Use = app.Config.RootCmd
}

// AppFromContext retrieves the App from the command context
func AppFromContext(ctx context.Context) *App {
	if a, ok := ctx.Value(appKey).(*App); ok {
		return a
	}
	// Fallback to package-level app (for testing or edge cases)
	return app
}

// ExecuteArgs executes the command with the given arguments using the fang library.
func ExecuteArgs(args []string) error {
	RootCmd.SetArgs(args)
	return fang.Execute(context.Background(), RootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme))
}

// ExecuteArgsString executes the command with the given arguments and returns the output as a string.
// This is useful for testing command outputs and used in the REPL
func ExecuteArgsString(args []string) (string, error) {
	return ExecuteArgsStringWithContext(context.Background(), args)
}

// ExecuteArgsStringWithContext executes the command with context and returns the output as a string.
func ExecuteArgsStringWithContext(ctx context.Context, args []string) (string, error) {
	buf := new(bytes.Buffer)

	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs(args)
	RootCmd.SetContext(ctx)

	err := RootCmd.Execute()

	return buf.String(), err
}
