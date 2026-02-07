package repl

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	words := strings.Fields(text)

	if len(words) == 0 {
		return nil
	}

	if words[0] != "pm" {
		return filterByPrefix(rootSuggestions, words[0])
	}

	pmWords := words[1:]
	if len(words) == 1 || (len(pmWords) == 1 && !strings.HasSuffix(text, " ")) {
		return commandSuggestions(pmWords)
	}

	if len(pmWords) >= 1 {
		return flagSuggestions(pmWords[0], pmWords, text)
	}

	return nil
}
