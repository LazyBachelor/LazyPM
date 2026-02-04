package dashboard

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type DashboardKeyMap struct {
	Quit        key.Binding
	SelectIssue key.Binding
	BackToList  key.Binding
	ScrollUp    key.Binding
	ScrollDown  key.Binding
}

var defaultDashboardKeyMap = DashboardKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	SelectIssue: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view issue"),
	),
	BackToList: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "back to list"),
	),
	ScrollUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "scroll up"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "scroll down"),
	),
}

func (d *Model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, d.keyMap.Quit):
		return tea.Quit
	case !d.focusedOnIssue && key.Matches(msg, d.keyMap.SelectIssue):
		d.focusedOnIssue = true
	case d.focusedOnIssue && key.Matches(msg, d.keyMap.BackToList):
		d.focusedOnIssue = false
	case d.focusedOnIssue && key.Matches(msg, d.keyMap.ScrollUp):
		d.issueView.Viewport.LineUp(1)
	case d.focusedOnIssue && key.Matches(msg, d.keyMap.ScrollDown):
		d.issueView.Viewport.LineDown(1)
	}

	return cmd
}
