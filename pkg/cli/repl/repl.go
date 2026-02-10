// Package repl implements the Read-Eval-Print Loop (REPL) for the PM CLI.
package repl

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli"
	"github.com/LazyBachelor/LazyPM/pkg/cli/commands"
	"github.com/LazyBachelor/LazyPM/pkg/cli/styles"
	"github.com/c-bata/go-prompt"
	"golang.org/x/term"
)

const (
	ReplHelp = `Type 'pm help' for available PM commands.
You can also run shell commands directly. Type 'exit' or 'quit' to leave.`

	ReplTitle = "Welcome to Project Management CLI! " + ReplHelp
)

// RunREPL starts the interactive Read-Eval-Print Loop for the PM CLI.
func RunREPL(ctx context.Context, config cli.CLIConfig) error {
	// Set terminal to raw mode to capture input properly in the REPL.
	// This allows us to handle input character by character and provide a better user experience.
	// We also ensure that the terminal state is restored when the REPL exits, even if an error occurs.
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to get terminal state: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Initialize services for beads, config and stats.
	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}
	defer cleanup()

	// Make sure to set services, to ensure they are available.
	commands.SetServices(svc)

	fmt.Println(styles.TitleStyle.Render(ReplTitle)) // Print REPL title.

	// history keeps track of command history.
	// This enables navigating through previous commands.
	var history []string

	// Start the REPL loop, which continues until the user types "exit" or "quit".
	for {
		// Prompt the user for input, and provide suggestions.
		input := prompt.Input(
			PromptPrefix,
			completer,
			promptOptions(history)...,
		)

		// Trim whitespace from the input to ensure consistent command processing.
		input = strings.TrimSpace(input)

		// If the user types "exit" or "quit", break the loop and exit the REPL.
		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Add the input to the history for future navigation.
		history = append(history, input)

		output, _ := execute(input)                     // Ignore errors for now, gives better ux
		fmt.Println(styles.CommandStyle.Render(output)) // Print the output of the command in a styled format.
	}

	return nil
}
