package repl

import "github.com/c-bata/go-prompt"

const PromptPrefix = "> "
const OptionMaxSuggestions = 5

// promptOptions returns a slice of prompt.Option
// to configure the behavior and appearance of the REPL prompt.
func promptOptions(history []string) []prompt.Option {
	return []prompt.Option{
		prompt.OptionPrefixTextColor(prompt.Cyan),
		prompt.OptionSuggestionTextColor(prompt.White),
		prompt.OptionMaxSuggestion(OptionMaxSuggestions),
		prompt.OptionSelectedSuggestionTextColor(prompt.Cyan),
		prompt.OptionSuggestionBGColor(prompt.DefaultColor),
		prompt.OptionSelectedSuggestionBGColor(prompt.DefaultColor),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionDescriptionBGColor(prompt.DefaultColor),
		prompt.OptionSelectedDescriptionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionBGColor(prompt.DefaultColor),
		prompt.OptionPreviewSuggestionBGColor(prompt.DefaultColor),
		prompt.OptionScrollbarBGColor(prompt.DefaultColor),
		prompt.OptionHistory(history),
	}
}
