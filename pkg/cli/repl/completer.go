package repl

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// completer provides suggestions for the REPL input based on the current input text.
func completer(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()  // Gets the text before the cursor as a string.
	words := strings.Fields(text) // Split the string into words as []string.

	// If there are no words, return no suggestions.
	if len(words) == 0 {
		return nil
	}

	// If the first word is not "pm", only provide root-level suggestions.
	if words[0] != "pm" {
		return filterByPrefix(rootSuggestions, words[0])
	}

	// If the first word is "pm", provide command and flag suggestions based on the context.
	pmWords := words[1:]
	if len(words) == 1 || (len(pmWords) == 1 && !strings.HasSuffix(text, " ")) {
		return commandSuggestions(pmWords)
	}

	// If the last word starts with a "-", provide flag suggestions for the current command.
	if len(pmWords) >= 1 {
		return flagSuggestions(pmWords[0], pmWords, text)
	}

	return nil
}
