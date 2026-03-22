// Package repl implements the Read-Eval-Print Loop (REPL) for the PM CLI.
package repl

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/app"
	issues "github.com/LazyBachelor/LazyPM/internal/commands/issues"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/style"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/c-bata/go-prompt"
	"golang.org/x/term"
)

type App = models.App
type ValidationFeedback = models.ValidationFeedback

const (
	ReplHelp = `Type 'pm help' for available PM commands.
You can also run shell commands directly. Type 'exit' or 'quit' to leave.
Type 'status' to check task progress.`

	ReplTitle = "Welcome to Project Management CLI! " + ReplHelp
)

type REPL struct {
	feedbackChan   chan ValidationFeedback
	quitChan       chan bool
	submitChan     chan<- struct{}
	completionChan chan struct{}
	app            *App

	currentFeedback ValidationFeedback
	exitRequested   bool
	taskCompleted   bool
}

func New() *REPL {
	return &REPL{}
}

// Run starts the interactive Read-Eval-Print Loop for the PM CLI.
func (r *REPL) Run(ctx context.Context, config app.Config) error {
	// Reset state for new task run
	r.taskCompleted = false
	r.exitRequested = false
	r.currentFeedback = ValidationFeedback{}
	r.completionChan = make(chan struct{}, 1)

	// Save terminal state to restore on exit
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err == nil {
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

	// Initialize services for beads, config and stats.
	app, cleanup, err := app.New(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}
	defer cleanup()

	// Make sure to set app, to ensure they are available.
	issues.SetApp(app)

	// Store app reference for updating feedback
	r.app = app

	if r.submitChan != nil {
		r.app.SubmitChan = r.submitChan
	}

	fmt.Println(style.TitleStyle.Render(ReplTitle))

	// Start goroutine to watch for validation feedback and quit signals
	if r.feedbackChan != nil && r.quitChan != nil {
		go r.watchValidation()
	}

	// Goroutine to inject newline when task completes to wake up blocked prompt
	go func() {
		for range r.completionChan {
			if tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0); err == nil {
				tty.Write([]byte("\n"))
				tty.Close()
			}
		}
	}()

	// history keeps track of command history.
	var history []string

	// Start the REPL loop
replLoop:
	for !r.exitRequested {
		// Check if we should exit before prompting (non-blocking check)
		if r.exitRequested {
			break
		}

		// Check if task completed before showing prompt
		if r.taskCompleted {
			fmt.Println("\nTask completed successfully!")
			break replLoop
		}

		// Prompt the user for input, and provide suggestions.
		input := prompt.Input(
			PromptPrefix,
			completer,
			promptOptions(history)...,
		)

		// Check if task completed while at prompt (will be true if newline was injected)
		if r.taskCompleted {
			fmt.Println("\nTask completed successfully!")
			break replLoop
		}

		// Check again after prompt returns (in case validation completed while waiting)
		if r.exitRequested {
			break
		}

		// Trim whitespace from the input
		input = strings.TrimSpace(input)

		// If the user types "exit" or "quit", break the loop and exit the REPL.
		if input == "exit" || input == "quit" {
			r.logAction("repl exit requested")
			fmt.Println("Goodbye!")
			break
		}

		if input != "" {
			r.logAction("repl command: " + input)
		}

		// Add the input to the history
		history = append(history, input)

		output, err := r.execute(input)

		// Send submit signal to trigger validation after any command
		if r.submitChan != nil {
			select {
			case r.submitChan <- struct{}{}:
			default:
			}
		}

		if err != nil {
			// Show command output (even on error)
			if output != "" {
				fmt.Println(style.TextStyle.Render(output))
			}
			// Show error message in red if no output
			if output == "" {
				fmt.Println(style.ErrorStyle.Render(err.Error()))
			}
		} else if output != "" {
			fmt.Println(style.TextStyle.Render(output))
		}
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
				r.app.CurrentFeedback = &ValidationFeedback{
					Success: feedback.Success,
					Message: feedback.Message,
					Checks:  feedback.Checks,
				}
			}
			if feedback.Success {
				r.taskCompleted = true
				if r.completionChan != nil {
					select {
					case r.completionChan <- struct{}{}:
					default:
					}
				}
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

func (r *REPL) SetSubmitChan(submitChan chan<- struct{}) {
	r.submitChan = submitChan
}

func (r *REPL) logAction(action string) {
	if r.app != nil {
		r.app.LogAction(models.EncodeActionEvent(models.ActionEvent{
			Source: "repl",
			Action: action,
		}))
	}
}
