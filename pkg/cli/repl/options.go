package repl

import "github.com/c-bata/go-prompt"

const PromptPrefix = "> "
const OptionMaxSuggestions = 5

const (
	ReplHelp = `Type 'pm help' for available PM commands.
You can also run shell commands directly. Type 'exit' or 'quit' to leave.`

	ReplTitle = "Welcome to Project Management CLI! " + ReplHelp
)

func promptOptions(history []string) []prompt.Option {
	return []prompt.Option{
		prompt.OptionPrefixTextColor(prompt.Cyan),
		prompt.OptionMaxSuggestion(OptionMaxSuggestions),
		prompt.OptionSuggestionBGColor(prompt.DefaultColor),
		prompt.OptionSelectedSuggestionBGColor(prompt.DefaultColor),
		prompt.OptionDescriptionBGColor(prompt.DefaultColor),
		prompt.OptionSelectedDescriptionBGColor(prompt.DefaultColor),
		prompt.OptionPreviewSuggestionBGColor(prompt.DefaultColor),
		prompt.OptionScrollbarBGColor(prompt.DefaultColor),
		prompt.OptionHistory(history),
	}
}
