package repl

import (
	"os/exec"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/commands/issues"
)

// execute processes the input command and returns the output
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

func executeShellCommand(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func executePMCommand(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}

	output, err := issuesCmd.ExecuteArgsString(parts)
	return output, err
}
