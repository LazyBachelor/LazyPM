package repl

import (
	"context"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/commands/issues"
	"github.com/c-bata/go-prompt"
	"github.com/muesli/reflow/truncate"
)

// rootSuggestions is a list of prompt suggestions for root-level commands.
var rootSuggestions = []prompt.Suggest{
	{Text: "pm", Description: "Project Management System"},
	{Text: "status", Description: "Show task status"},
	{Text: "exit", Description: "Exit pm CLI"},
	{Text: "help", Description: "Show help information"},
	{Text: "title", Description: "Print the welcome title"},
	{Text: "git", Description: "Version control system"},
}

// commandSuggestions is a list of prompt suggestions for PM commands.
var baseSuggestions = []prompt.Suggest{
	{Text: "help", Description: "Show help information"},
	{Text: "sprint", Description: "Manage sprints"},
	{Text: "dep", Description: "Manage dependencies"},
	{Text: "create", Description: "Create a new issue with title"},
	{Text: "list", Description: "List all issues"},
	{Text: "read", Description: "Read issue details by ID"},
	{Text: "update", Description: "Update an existing issue by ID"},
	{Text: "close", Description: "Close an issue by ID"},
	{Text: "delete", Description: "Delete an issue by ID"},
	{Text: "comment", Description: "Add a comment on an issue by ID"},
	{Text: "comments", Description: "List comments on an issue by ID"},

	{Text: "new", Description: "Alias for create command"},
	{Text: "ls", Description: "Alias for list command"},
	{Text: "search", Description: "Alias for list command"},
	{Text: "describe", Description: "Alias for read command"},
	{Text: "details", Description: "Alias for read command"},
	{Text: "edit", Description: "Alias for update command"},
	{Text: "rm", Description: "Alias for delete command"},
	{Text: "del", Description: "Alias for delete command"},
	{Text: "get", Description: "Alias for read command"},
	{Text: "dependencies", Description: "Alias for dep"},
}

// createFlags is a list of prompt suggestions for the create command flags.
var createFlags = []prompt.Suggest{
	{Text: "--desc", Description: "Issue description"},
	{Text: "--status", Description: "Issue status"},
	{Text: "--type", Description: "Issue type"},
	{Text: "--priority", Description: "Issue priority"},
	{Text: "--assignee", Description: "Issue assignee"},

	{Text: "-d", Description: "Issue description (short)"},
	{Text: "-s", Description: "Issue status (short)"},
	{Text: "-t", Description: "Issue type (short)"},
	{Text: "-p", Description: "Issue priority (short)"},
	{Text: "-a", Description: "Issue assignee (short)"},
}

// updateFlags is a list of prompt suggestions for the update command flags.
var updateFlags = []prompt.Suggest{
	{Text: "--title", Description: "New issue title"},
	{Text: "--desc", Description: "New issue description"},
	{Text: "--status", Description: "New issue status"},
	{Text: "--type", Description: "New issue type"},
	{Text: "--priority", Description: "New issue priority"},
	{Text: "--assignee", Description: "New issue assignee"},

	{Text: "-t", Description: "New issue type (short)"},
	{Text: "-p", Description: "New issue priority (short)"},
	{Text: "-a", Description: "New issue assignee (short)"},
	{Text: "-d", Description: "New issue description (short)"},
	{Text: "-s", Description: "New issue status (short)"},
}

// listFlags is a list of prompt suggestions for the list command flags.
var listFlags = []prompt.Suggest{
	{Text: "--title", Description: "Filter by title"},
	{Text: "--desc", Description: "Filter by description"},
	{Text: "--status", Description: "Filter by status"},
	{Text: "--type", Description: "Filter by type"},
	{Text: "--priority", Description: "Filter by priority"},
	{Text: "--assignee", Description: "Filter by assignee"},
	{Text: "--limit", Description: "Limit number of results"},

	{Text: "-s", Description: "Filter by status (short)"},
	{Text: "-t", Description: "Filter by type (short)"},
	{Text: "-p", Description: "Filter by priority (short)"},
	{Text: "-a", Description: "Filter by assignee (short)"},
	{Text: "-d", Description: "Filter by description (short)"},
	{Text: "-l", Description: "Limit number of results (short)"},
}

var deleteFlags = []prompt.Suggest{
	{Text: "--yes", Description: "Confirm deletion without prompt"},
	{Text: "--interactive", Description: "Select issues to delete interactively"},
}

var commentFlags = []prompt.Suggest{
	{Text: "--message", Description: "Comment text"},
	{Text: "--author", Description: "Author name for the comment"},

	{Text: "-m", Description: "Comment text (short)"},
	{Text: "-a", Description: "Author name (short)"},
}

// sprintSubcommands is a list of prompt suggestions for sprint subcommands.
var sprintSubcommands = []prompt.Suggest{
	{Text: "list", Description: "List all sprints"},
	{Text: "create", Description: "Create a new sprint"},
	{Text: "issues", Description: "List issues in a sprint"},
	{Text: "add", Description: "Add an issue to a sprint"},
	{Text: "remove", Description: "Remove an issue from a sprint"},
	{Text: "backlog", Description: "Show backlog issues"},
	{Text: "delete", Description: "Delete a sprint"},

	{Text: "ls", Description: "Alias for list subcommand"},
	{Text: "new", Description: "Alias for create subcommand"},
	{Text: "show", Description: "Alias for issues subcommand"},
	{Text: "view", Description: "Alias for issues subcommand"},
	{Text: "rm", Description: "Alias for remove subcommand"},
	{Text: "del", Description: "Alias for delete subcommand"},
}

var depSubcommands = []prompt.Suggest{
	{Text: "view", Description: "View dependencies of an issue"},
	{Text: "show", Description: "Alias for view"},
	{Text: "add", Description: "Add a dependency"},
	{Text: "remove", Description: "Remove a dependency"},
	{Text: "rm", Description: "Alias for remove"},
}

// statusValues is a list of prompt suggestions for status types
var statusValues = []prompt.Suggest{
	{Text: "open", Description: "Open status"},
	{Text: "closed", Description: "Closed status"},
	{Text: "in_progress", Description: "In progress status"},
	{Text: "blocked", Description: "Blocked status"},
}

// typeValues is a list of prompt suggestions for issue types
var typeValues = []prompt.Suggest{
	{Text: "task", Description: "Task issue type"},
	{Text: "bug", Description: "Bug issue type"},
	{Text: "feature", Description: "Feature issue type"},
	{Text: "chore", Description: "Chore issue type"},
}

// priorityValues is a list of prompt suggestions for issue priority levels
var priorityValues = []prompt.Suggest{
	{Text: "0", Description: "Irrelevant"},
	{Text: "1", Description: "Low priority"},
	{Text: "2", Description: "Normal priority"},
	{Text: "3", Description: "High priority"},
	{Text: "4", Description: "Critical priority"},
}

// isIDCommand maps command names to a boolean indicating whether they expect an issue ID as an argument.
var isIDCommand = map[string]bool{
	"describe":      true,
	"delete":        true,
	"del":           true,
	"rm":            true,
	"remove":        true,
	"get":           true,
	"read":          true,
	"close":         true,
	"update":        true,
	"edit":          true,
	"comment":       true,
	"comments":      true,
	"sprint-add":    true,
	"sprint-remove": true,
}

var commandFlags = map[string][]prompt.Suggest{
	"create":   createFlags,
	"add":      createFlags,
	"update":   updateFlags,
	"edit":     updateFlags,
	"list":     listFlags,
	"ls":       listFlags,
	"search":   listFlags,
	"delete":   deleteFlags,
	"del":      deleteFlags,
	"rm":       deleteFlags,
	"remove":   deleteFlags,
	"comment":  commentFlags,
	"comments": nil, // no flags, just issue ID
	"sprint":   sprintSubcommands,
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

	// Custom completion logic for `pm dep ...`
	if cmd == "dep" {
		switch len(words) {
		case 1:
			// `pm dep ` => next token could be a subcommand; default view is handled
			// once the user types an issue id
			if strings.HasSuffix(text, " ") {
				return depSubcommands
			}
			return nil

		case 2:
			// either:
			// - pm dep <issueID> (default view)
			// - pm dep <subcommand>
			if strings.HasSuffix(text, " ") {
				sub := words[1]
				switch sub {
				case "view", "show", "add", "remove", "rm":
					// After a valid subcommand we expect the issue id next.
					return issueIDSuggestions("", true)
				default:
					// Treat it as an issue id already; no more suggestions.
					return nil
				}
			}

			// Typing the 2nd token: suggest subcommands and issue IDs that match.
			suggests := make([]prompt.Suggest, 0, len(depSubcommands))
			suggests = append(suggests, filterByPrefix(depSubcommands, lastWord)...)
			//suggests = append(suggests, issueIDSuggestions(lastWord, true)...)
			return suggests

		case 3:
			// If subcommand is `add/remove`, third token is the issue id.
			// If subcommand is show/list, third token is the issue id.
			sub := words[1]
			switch sub {
			case "add", "remove", "rm":
				if strings.HasSuffix(text, " ") {
					// `pm dep add <issueID> ` => next token is dependsOnID.
					return issueIDSuggestions("", true)
				}
				return issueIDSuggestions(lastWord, true)
			default:
				// show/list/show => completing the issue id.
				if strings.HasSuffix(text, " ") {
					return nil
				}
				return issueIDSuggestions(lastWord, true)
			}

		case 4:
			// add/remove/rm have a 4th token: dependsOnID.
			sub := words[1]
			if sub == "add" || sub == "remove" || sub == "rm" {
				if strings.HasSuffix(text, " ") {
					return nil
				}
				return issueIDSuggestions(lastWord, true)
			}
			return nil
		}

		return nil
	}

	if values := getFlagValues(prevWord); values != nil {
		return filterByPrefix(values, lastWord)
	}

	if cmd == "sprint" {
		switch len(words) {
		case 1:
			if strings.HasSuffix(text, " ") {
				return sprintSubcommands
			}
			return nil
		case 2:
			subcommand := words[1]
			if !strings.HasSuffix(text, " ") {
				return filterByPrefix(sprintSubcommands, lastWord)
			}

			switch subcommand {
			case "list", "ls", "create", "new", "backlog":
				return nil
			case "issues", "show", "view", "delete", "del":
				return sprintNumSuggestions("", subcommand != "delete" && subcommand != "del")
			case "add", "remove", "rm":
				return issueIDSuggestions("", true)
			default:
				return filterByPrefix(sprintSubcommands, subcommand)
			}
		case 3:
			subcommand := words[1]

			switch subcommand {
			case "issues", "show", "view", "delete", "del":
				if !strings.HasSuffix(text, " ") {
					return sprintNumSuggestions(lastWord, subcommand != "delete" && subcommand != "del")
				}
				return nil
			case "add", "remove", "rm":
				if !strings.HasSuffix(text, " ") {
					return issueIDSuggestions(lastWord, true)
				}
				return sprintNumSuggestions("", true)
			}
		case 4:
			subcommand := words[1]

			switch subcommand {
			case "add", "remove", "rm":
				if !strings.HasSuffix(text, " ") {
					return sprintNumSuggestions(lastWord, true)
				}
			}
		}
		return nil
	}

	flags := commandFlags[cmd]

	if isIDCommand[cmd] {
		if len(words) == 1 && strings.HasSuffix(text, " ") {
			return issueIDSuggestions("", true)
		}

		if len(words) == 2 &&
			!strings.HasSuffix(text, " ") &&
			!strings.HasPrefix(words[1], "-") {
			return issueIDSuggestions(words[1], true)
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

	issues, _ := issues.GetIssueCompletions(context.Background(), partial)

	var suggestions []prompt.Suggest
	for _, issue := range issues {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        issue.ID,
			Description: truncate.String(issue.Title, 20),
		})
	}
	return suggestions
}

// sprintNumSuggestions returns a list of prompt suggestions for sprint numbers.
func sprintNumSuggestions(partial string, includeBacklog bool) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "0", Description: "Backlog"},
		{Text: "1", Description: "Sprint 1"},
		{Text: "2", Description: "Sprint 2"},
		{Text: "3", Description: "Sprint 3"},
		{Text: "4", Description: "Sprint 4"},
		{Text: "5", Description: "Sprint 5"},
	}

	if !includeBacklog {
		suggestions = suggestions[1:]
	}

	return filterByPrefix(suggestions, partial)
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
