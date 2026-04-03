package repl

import (
	"os/exec"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/commands/issues"
)

// execute processes the input command and returns the output
func (r *REPL) execute(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	if input == "help" {
		return ReplHelp, nil
	}

	if input == "title" {
		return ReplTitle, nil
	}

	if input == "status" {
		return executePMCommand("status")
	}

	if input == "pm" {
		return executePMCommand("help")
	}

	if after, ok := strings.CutPrefix(input, "pm"); ok {
		return executePMCommand(after)
	}
	return executeShellCommand(input)
}

// shellSplit splits input respecting quoted strings and escapes
// similar to how a shell would parse arguments
func shellSplit(input string) []string {
	var args []string
	var current strings.Builder
	var inQuote rune
	var escaped bool

	for _, ch := range input {
		if escaped {
			current.WriteRune(ch)
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if inQuote != 0 {
			if ch == inQuote {
				inQuote = 0
				continue
			}
			current.WriteRune(ch)
			continue
		}

		if ch == '"' || ch == '\'' {
			inQuote = ch
			continue
		}

		if ch == ' ' || ch == '\t' {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteRune(ch)
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

func executeShellCommand(input string) (string, error) {
	parts := shellSplit(input)
	if len(parts) == 0 {
		return "", nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func executePMCommand(input string) (string, error) {
	parts := shellSplit(input)
	if len(parts) == 0 {
		return "", nil
	}

	output, err := issues.ExecuteArgsString(parts)
	return output, err
}
