// Package repl implements the Read-Eval-Print Loop (REPL) for the PM CLI.
package repl

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/commands/issues"
	surveyCmd "github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/internal/style"
	"github.com/LazyBachelor/LazyPM/pkg/cli"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/c-bata/go-prompt"
	"golang.org/x/term"
)

const (
	ReplHelp = `Type 'pm help' for available PM commands.
Type 'pm status' to check task progress.
You can also run shell commands directly. Type 'exit' or 'quit' to leave.`

	ReplTitle = "Welcome to Project Management CLI! " + ReplHelp
)

type REPL struct {
	feedbackChan chan task.ValidationFeedback
	quitChan     chan bool
	app          *service.App

	currentFeedback task.ValidationFeedback
	exitRequested   bool
}

func NewRepl() *REPL {
	return &REPL{}
}

// Run starts the interactive Read-Eval-Print Loop for the PM CLI.
func (r *REPL) Run(ctx context.Context, config cli.Config) error {
	// Set terminal to raw mode to capture input properly in the REPL.
	// This allows us to handle input character by character and provide a better user experience.
	// We also ensure that the terminal state is restored when the REPL exits, even if an error occurs.
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to get terminal state: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Initialize services for beads, config and stats.
	app, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}
	defer cleanup()

	// Make sure to set app, to ensure they are available.
	issuesCmd.SetApp(app)

	// Store app reference for updating feedback
	r.app = app

	fmt.Println(style.TitleStyle.Render(ReplTitle)) // Print REPL title.

	// Start goroutine to watch for validation feedback and quit signals
	if r.feedbackChan != nil && r.quitChan != nil {
		go r.watchValidation()
	}

	// history keeps track of command history.
	// This enables navigating through previous commands.
	var history []string

	// Start the REPL loop, which continues until the user types "exit" or "quit" or task completes.
	for !r.exitRequested {
		// Check if we should exit before prompting (non-blocking check)
		if r.exitRequested {
			break
		}

		// Prompt the user for input, and provide suggestions.
		input := prompt.Input(
			PromptPrefix,
			completer,
			promptOptions(history)...,
		)

		// Check again after prompt returns (in case validation completed while waiting)
		if r.exitRequested {
			break
		}

		// Trim whitespace from the input to ensure consistent command processing.
		input = strings.TrimSpace(input)

		// If the user types "exit" or "quit", break the loop and exit the REPL.
		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Add the input to the history for future navigation.
		history = append(history, input)

		output, _ := execute(input)                 // Ignore errors for now, gives better ux
		fmt.Println(style.TextStyle.Render(output)) // Print the output of the command in a styled format.
	}

	return nil
}

func (r *REPL) watchValidation() {
	for {
		select {
		case feedback := <-r.feedbackChan:
			r.currentFeedback = feedback
			// Update app's CurrentFeedback so status command can access it
			if r.app != nil {
				r.app.CurrentFeedback = &service.ValidationFeedback{
					Success: feedback.Success,
					Message: feedback.Message,
				}
			}
			if feedback.Success {
				fmt.Printf("\n%s\n", styles.TitleStyle.Render("Task completed successfully!"))
				fmt.Println("Press Enter to exit...")
				r.exitRequested = true
				return
			}
		case <-r.quitChan:
			r.exitRequested = true
			return
		}
	}
}

// SetChannels sets the channels for receiving validation feedback and quit signals from the task interface
func (r *REPL) SetChannels(feedbackChan chan task.ValidationFeedback, quitChan chan bool) {
	r.feedbackChan = feedbackChan
	r.quitChan = quitChan
}

func init() {
	issuesCmd.RootCmd.AddCommand(issuesCmd.GetCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.ListCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.CloseCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.CreateCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.DeleteCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.UpdateCmd)
	issuesCmd.RootCmd.AddCommand(surveyCmd.StatusCmd)
}
