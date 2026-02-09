package repl

import (
	"context"
	"strings"

	"github.com/LazyBachelor/LazyPM/pkg/cli/commands"
	"github.com/c-bata/go-prompt"
	"github.com/muesli/reflow/truncate"
)

// rootSuggestions is a list of prompt suggestions for root-level commands.
var rootSuggestions = []prompt.Suggest{
	{Text: "pm", Description: "Project Management System"},
	{Text: "exit", Description: "Exit pm CLI"},
	{Text: "help", Description: "Show help information"},
	{Text: "title", Description: "Print the welcome title"},
	{Text: "git", Description: "Version control system"},
}

// commandSuggestions is a list of prompt suggestions for PM commands.
var baseSuggestions = []prompt.Suggest{
	{Text: "help", Description: "Show help information"},
	{Text: "delete", Description: "Delete an issue by ID"},
	{Text: "close", Description: "Close an issue by ID"},
	{Text: "create", Description: "Create a new issue with title"},
	{Text: "update", Description: "Update an existing issue by ID"},
	{Text: "update", Description: "Update an existing issue by ID"},
	{Text: "describe", Description: "Get issue details by ID"},
	{Text: "list", Description: "List all issues"},
}

// createFlags is a list of prompt suggestions for the create command flags.
var createFlags = []prompt.Suggest{
	{Text: "--desc", Description: "Issue description"},
	{Text: "--status", Description: "Issue status (open, closed, in_progress)"},
	{Text: "--type", Description: "Issue type (bug, feature, task)"},
	{Text: "--priority", Description: "Issue priority (0-5)"},
}

// updateFlags is a list of prompt suggestions for the update command flags.
var updateFlags = []prompt.Suggest{
	{Text: "--title", Description: "New issue title"},
	{Text: "--desc", Description: "New issue description"},
	{Text: "--status", Description: "New issue status (open, closed, in_progress)"},
	{Text: "--type", Description: "New issue type (bug, feature, task)"},
	{Text: "--priority", Description: "New issue priority (0-5)"},
}

// listFlags is a list of prompt suggestions for the list command flags.
var listFlags = []prompt.Suggest{
	{Text: "--title", Description: "Filter by title"},
	{Text: "--desc", Description: "Filter by description"},
	{Text: "--status", Description: "Filter by status (open, closed, in_progress)"},
	{Text: "--type", Description: "Filter by type (bug, feature, task)"},
	{Text: "--priority", Description: "Filter by priority (0-5)"},
	{Text: "--limit", Description: "Limit number of results"},
}

// statusValues is a list of prompt suggestions for status types
var statusValues = []prompt.Suggest{
	{Text: "open", Description: "Open status"},
	{Text: "closed", Description: "Closed status"},
	{Text: "in_progress", Description: "In progress status"},
}

// typeValues is a list of prompt suggestions for issue types
var typeValues = []prompt.Suggest{
	{Text: "bug", Description: "Bug issue type"},
	{Text: "feature", Description: "Feature issue type"},
	{Text: "task", Description: "Task issue type"},
}

// priorityValues is a list of prompt suggestions for issue priority levels
var priorityValues = []prompt.Suggest{
	{Text: "0", Description: "Lowest priority"},
	{Text: "1", Description: "Low priority"},
	{Text: "2", Description: "Medium-low priority"},
	{Text: "3", Description: "Medium priority"},
	{Text: "4", Description: "High priority"},
}

// isIDCommand maps command names to a boolean indicating whether they expect an issue ID as an argument.
var isIDCommand = map[string]bool{
	"describe": true,
	"delete":   true,
	"del":      true,
	"rm":       true,
	"remove":   true,
	"get":      true,
	"read":     true,
	"close":    true,
	"update":   true,
	"edit":     true,
}

var commandFlags = map[string][]prompt.Suggest{
	"create": createFlags,
	"add":    createFlags,
	"update": updateFlags,
	"edit":   updateFlags,
	"list":   listFlags,
	"ls":     listFlags,
	"search": listFlags,
}

// commandSuggestions returns a list of prompt suggestions based on the current input words.
func commandSuggestions(words []string) []prompt.Suggest {
	if len(words) == 0 {
		return baseSuggestions
	}
	return filterByPrefix(baseSuggestions, words[0])
}

// flagSuggestions returns a list of prompt suggestions for command flags based on the current input.
func flagSuggestions(cmd string, words []string, text string) []prompt.Suggest {
	lastWord, prevWord := parseWords(words, text)

	if values := getFlagValues(prevWord); values != nil {
		return filterByPrefix(values, lastWord)
	}

	flags := commandFlags[cmd]

	if isIDCommand[cmd] {
		if len(words) < 2 && !strings.HasPrefix(lastWord, "-") {
			return issueIDSuggestions(lastWord, true)
		}
		return filterByPrefix(flags, lastWord)
	}

	return filterByPrefix(flags, lastWord)
}

// issueIDSuggestions returns a list of prompt suggestions for issue IDs based on the current partial input.
func issueIDSuggestions(partial string, hasCommand bool) []prompt.Suggest {
	// Only show suggestions if we've typed the command already
	if !hasCommand {
		return nil
	}

	issues, _ := commands.GetIssueCompletions(context.Background(), partial)

	var suggestions []prompt.Suggest
	for _, issue := range issues {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        issue.ID,
			Description: truncate.String(issue.Title, 20),
		})
	}
	return suggestions
}

// parseWords extracts the last and previous words from the input for flag suggestion logic.
func parseWords(words []string, text string) (lastWord, prevWord string) {
	if len(words) > 0 && !strings.HasSuffix(text, " ") {
		lastWord = words[len(words)-1]
		if len(words) >= 2 {
			prevWord = words[len(words)-2]
		}
	} else if len(words) >= 1 {
		prevWord = words[len(words)-1]
	}
	return
}

// getFlagValues returns a list of prompt suggestions for flag values based on the given flag.
func getFlagValues(flag string) []prompt.Suggest {
	switch flag {
	case "-s", "--status":
		return statusValues
	case "-t", "--type":
		return typeValues
	case "-p", "--priority":
		return priorityValues
	}
	return nil
}

// filterByPrefix filters a list of prompt suggestions based on a given prefix.
func filterByPrefix(suggestions []prompt.Suggest, prefix string) []prompt.Suggest {
	if prefix == "" {
		return suggestions
	}
	var filtered []prompt.Suggest
	for _, s := range suggestions {
		if strings.HasPrefix(s.Text, prefix) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
