package repl

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/cli"
	"github.com/LazyBachelor/LazyPM/pkg/cli/commands"
	"github.com/LazyBachelor/LazyPM/pkg/cli/styles"
	"github.com/c-bata/go-prompt"
	"github.com/muesli/reflow/truncate"
	"golang.org/x/term"
)

const (
	ReplHelp = `Type 'pm help' for available PM commands.
You can also run shell commands directly. Type 'exit' or 'quit' to leave.`

	ReplTitle = "Welcome to Project Management CLI! " + ReplHelp
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

	fmt.Println("\n" + styles.TitleStyle.Render(ReplTitle))

	var history []string

	for {
		input := prompt.Input(
			"› ",
			completer,
			prompt.OptionPrefixTextColor(prompt.Cyan),
			prompt.OptionMaxSuggestion(5),
			prompt.OptionSuggestionBGColor(prompt.DefaultColor),
			prompt.OptionSelectedSuggestionBGColor(prompt.DefaultColor),
			prompt.OptionDescriptionBGColor(prompt.DefaultColor),
			prompt.OptionSelectedDescriptionBGColor(prompt.DefaultColor),
			prompt.OptionPreviewSuggestionBGColor(prompt.DefaultColor),
			prompt.OptionScrollbarBGColor(prompt.DefaultColor),
			prompt.OptionHistory(history),
		)

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if input == "help" {
			fmt.Println(ReplHelp)
			continue
		}

		if input == "title" {
			fmt.Println(styles.TitleStyle.Render(ReplTitle))
			continue
		}

		history = append(history, input)

		if after, ok := strings.CutPrefix(input, "pm "); ok {
			pmCmd := strings.TrimSpace(after)
			if pmCmd == "" {
				continue
			}

			args := strings.Fields(pmCmd)

			output, _ := commands.ExecuteArgsString(args)

			if output != "" {
				fmt.Println(styles.CommandStyle.Render(output))
			}

		} else {
			cmd := exec.Command("sh", "-c", input)
			out, err := cmd.CombinedOutput()

			if len(out) > 0 {
				fmt.Println(styles.CommandStyle.Render(string(out)))
			}

			if err != nil && len(out) == 0 {
				fmt.Fprintf(os.Stderr, "%s", fmt.Sprintf("Error: %v", err))
			}
		}
	}

	return nil
}

func completer(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	words := strings.Fields(text)

	if len(words) == 0 {
		return nil
	}

	if words[0] != "pm" {
		suggestions := []prompt.Suggest{
			{Text: "pm", Description: "Project Management System"},
			{Text: "exit", Description: "Exit pm CLI"},
			{Text: "help", Description: "Show help information"},
			{Text: "title", Description: "Print the welcome title"},
			{Text: "git", Description: "Version control system"},
		}

		if len(words) == 0 {
			return suggestions
		}

		var filtered []prompt.Suggest
		for _, s := range suggestions {
			if strings.HasPrefix(s.Text, words[0]) {
				filtered = append(filtered, s)
			}
		}

		return filtered
	}

	// If first word is "pm", provide PM command completions
	if words[0] == "pm" {
		// Remove "pm" from words to get the actual command
		if len(words) == 1 || (len(words) == 2 && !strings.HasSuffix(text, " ")) {
			return commandSuggestions(words[1:])
		}

		// Get the PM subcommand
		if len(words) >= 2 {
			cmd := words[1]
			return flagSuggestions(cmd, words[1:], text)
		}
	}

	return nil
}

func commandSuggestions(words []string) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "help", Description: "Show help information"},
		{Text: "delete", Description: "Delete an issue by ID"},
		{Text: "create", Description: "Create a new issue with title"},
		{Text: "describe", Description: "Get issue details by ID"},
		{Text: "list", Description: "List all issues"},
	}

	if len(words) == 0 {
		return suggestions
	}

	var filtered []prompt.Suggest
	for _, s := range suggestions {
		if strings.HasPrefix(s.Text, words[0]) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func flagSuggestions(cmd string, words []string, text string) []prompt.Suggest {
	var flagSuggests []prompt.Suggest
	isCompleteWord := strings.HasSuffix(text, " ")
	lastWord := ""
	if len(words) > 0 && !isCompleteWord {
		lastWord = words[len(words)-1]
	}

	switch cmd {
	case "create", "add":
		flagSuggests = []prompt.Suggest{
			{Text: "--desc", Description: "Issue description"},
			{Text: "--status", Description: "Issue status (open, closed, in_progress)"},
			{Text: "--type", Description: "Issue type (bug, feature, task)"},
			{Text: "--priority", Description: "Issue priority (0-5)"},
		}
		return filterAndCompleteFlags(flagSuggests, lastWord, words)

	case "ls", "list", "search":
		flagSuggests = []prompt.Suggest{
			{Text: "--title", Description: "Filter by title"},
			{Text: "--desc", Description: "Filter by description"},
			{Text: "--status", Description: "Filter by status (open, closed, in_progress)"},
			{Text: "--type", Description: "Filter by type (bug, feature, task)"},
			{Text: "--priority", Description: "Filter by priority (0-5)"},
			{Text: "--limit", Description: "Limit number of results"},
		}
		return filterAndCompleteFlags(flagSuggests, lastWord, words)

	case "describe", "get", "read":
		return issueIdSuggestions(words)

	case "delete", "del", "rm", "remove":
		return issueIdSuggestions(words)
	}

	return nil
}

func issueIdSuggestions(words []string) []prompt.Suggest {
	if len(words) < 2 {
		return nil
	}

	// Get the partial ID being typed
	partial := ""
	if len(words) >= 2 {
		partial = words[len(words)-1]
	}

	issues, _ := commands.GetIssueCompletions(context.Background(), partial)

	var suggestions []prompt.Suggest
	for _, issue := range issues {
		suggestions = append(suggestions, prompt.Suggest{
			Text: issue.ID, Description: truncate.String(issue.Title, 20),
		})
	}
	return suggestions
}

func filterAndCompleteFlags(suggestions []prompt.Suggest, lastWord string, words []string) []prompt.Suggest {
	if len(words) >= 2 {
		prevWord := words[len(words)-2]
		switch prevWord {
		case "-s", "--status":
			return []prompt.Suggest{
				{Text: "open", Description: "Open status"},
				{Text: "closed", Description: "Closed status"},
				{Text: "in_progress", Description: "In progress status"},
			}
		case "-t", "--type":
			return []prompt.Suggest{
				{Text: "bug", Description: "Bug issue type"},
				{Text: "feature", Description: "Feature issue type"},
				{Text: "task", Description: "Task issue type"},
			}
		case "-p", "--priority":
			return []prompt.Suggest{
				{Text: "0", Description: "Lowest priority"},
				{Text: "1", Description: "Low priority"},
				{Text: "2", Description: "Medium-low priority"},
				{Text: "3", Description: "Medium priority"},
				{Text: "4", Description: "High priority"},
				{Text: "5", Description: "Highest priority"},
			}
		}
	}

	if lastWord == "" {
		return suggestions
	}

	var filtered []prompt.Suggest
	for _, s := range suggestions {
		if strings.HasPrefix(s.Text, lastWord) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
