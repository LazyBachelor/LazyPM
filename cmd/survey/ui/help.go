package ui

import "github.com/charmbracelet/bubbles/key"

type TaskHelpKeys struct {
	Quit     key.Binding
	Continue key.Binding
}

var DefaultTaskKeys = TaskHelpKeys{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "Quit"),
	),
	Continue: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "Continue"),
	),
}

func (h TaskHelpKeys) ShortHelp() []key.Binding {
	return []key.Binding{h.Continue, h.Quit}
}

func (h TaskHelpKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{{h.Continue, h.Quit}}
}
