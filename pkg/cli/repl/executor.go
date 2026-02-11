package repl

import (
	"os/exec"
	"strings"

	"github.com/LazyBachelor/LazyPM/pkg/cli/commands"
)

// execute processes the input command and returns the output
// or an error if it occurs. It handles what type of command is being executed,
// whether it's a PM command or a shell command, and routes it accordingly.
func execute(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	if input == "help" {
		return ReplHelp, nil
	}

	if input == "title" {
		return ReplTitle, nil
	}

	if after, ok := strings.CutPrefix(input, "pm"); ok {
		return executePMCommand(after)
	}
	return executeShellCommand(input)
}

// executeShellCommand executes a shell command
// and returns its output or an error if it occurs.
func executeShellCommand(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// executePMCommand executes a PM command using the commands package
// and returns its output or an error if it occurs.
func executePMCommand(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}

	output, err := commands.ExecuteArgsString(parts)
	return output, err
}
