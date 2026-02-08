package commands

import (
	"bytes"
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/charmbracelet/fang"

	"github.com/spf13/cobra"
)

// svc is a global variable that holds beads, config and stats services.
// Must be called before executing any commands to ensure services are available.
var svc *service.Services

type Flags struct {
	interactive bool
	title       string
	description string
	status      string
	issueType   string
	priority    int
}

// rootCmd is the base command for the CLI application.
var rootCmd = &cobra.Command{
	Short: "Project Management CLI",
	Long:  `Project Management CLI for managing issues and tasks.`,
}

// SetServices sets the global services variable for use in command execution.
// Must be called before executing any commands to ensure services are available.
func SetServices(services *service.Services) {
	svc = services
	rootCmd.Use = svc.Config.RootCmd
}

// Execute executes the root command using the fang library.
func Execute() error {
	return fang.Execute(context.Background(), rootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme))
}

// ExecuteArgs executes the command with the given arguments using the fang library.
func ExecuteArgs(args []string) error {
	rootCmd.SetArgs(args)
	return fang.Execute(context.Background(), rootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme))
}

// ExecuteArgsString executes the command with the given arguments and returns the output as a string.
// This is useful for testing command outputs and used in the REPL
func ExecuteArgsString(args []string) (string, error) {
	buf := new(bytes.Buffer)

	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()

	return buf.String(), err
}

// init function to set up the command hierarchy and options.
func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.AddGroup(&cobra.Group{ID: "help", Title: "Helping Commands"})
	rootCmd.SetCompletionCommandGroupID("help")
	rootCmd.SetHelpCommandGroupID("help")
}
