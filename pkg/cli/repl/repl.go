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

func RunREPL(ctx context.Context, config cli.CLIConfig) error {
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to get terminal state: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}
	defer cleanup()

	commands.SetServices(svc)

	fmt.Println(styles.TitleStyle.Render(ReplTitle))

	var history []string

	for {
		input := prompt.Input(
			PromptPrefix,
			completer,
			promptOptions(history)...,
		)

		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		history = append(history, input)

		// Ignore errors for now, gives better ux
		output, _ := execute(input)

		fmt.Println(styles.CommandStyle.Render(output))
	}

	return nil
}
