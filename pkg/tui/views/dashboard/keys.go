package dashboard

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
)

type KeyMap struct {
	components.CommonKeyMap
	SwitchToKanbanBoard key.Binding
	Quit                key.Binding
	SelectIssue         key.Binding
	BackToList          key.Binding
	ScrollUp            key.Binding
	ScrollDown          key.Binding
	EditTitle           key.Binding
	EditDescription     key.Binding
	ChangeStatus        key.Binding
	ChangePriority      key.Binding
	ChangeType          key.Binding
	ChangeAssignee      key.Binding
	AddComment          key.Binding
	AddIssue            key.Binding
	DeleteIssue         key.Binding
}

var defaultDashboardKeyMap = KeyMap{
	CommonKeyMap: components.DefaultCommonKeyMap(),
	SwitchToKanbanBoard: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "switch to kanban")),
	EditTitle: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit title"),
	),
	EditDescription: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "edit description"),
	),
	ChangeStatus: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "change status"),
	),
	ChangePriority: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "change priority"),
	),
	ChangeType: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "change type"),
	),
	ChangeAssignee: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "change assignee"),
	),
	AddComment: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "add comment"),
	),
	AddIssue: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add issue"),
	),
	DeleteIssue: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete issue"),
	),
}

func (m *Model) handleKeyPressMsg(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keyMap.Help):
		m.helpBar.ToggleHelp()
		m.logAction("tui toggled help")

	case key.Matches(msg, m.keyMap.Quit):
		m.logAction("tui quit requested")
		return tea.Quit

	case !m.IsInModal() && key.Matches(msg, m.keyMap.SwitchToKanbanBoard):
		return func() tea.Msg { return msgs.SwitchToKanbanBoardMsg{} }

	case m.IsFocusedOnDetail() && key.Matches(msg, m.keyMap.ScrollUp):
		m.issueDetail.ScrollUp(1)
		m.logAction("tui scrolled issue detail up")

	case m.IsFocusedOnDetail() && key.Matches(msg, m.keyMap.ScrollDown):
		m.issueDetail.ScrollDown(1)
		m.logAction("tui scrolled issue detail down")

	case !m.IsInModal() && key.Matches(msg, m.keyMap.EditTitle):
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			cmd = m.startEditTitle(selected)
			m.logAction("tui started editing issue title")
		}

	case !m.IsInModal() && key.Matches(msg, m.keyMap.EditDescription):
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			cmd = m.startEditDescription(selected)
			m.logAction("tui started editing issue description")
		}

	case !m.IsInModal() && key.Matches(msg, m.keyMap.ChangeStatus):
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			cmd = m.startChooseStatus(selected)
			m.logAction("tui opened status picker")
		}

	case !m.IsInModal() && key.Matches(msg, m.keyMap.ChangePriority):
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			cmd = m.startChoosePriority(selected)
			m.logAction("tui opened priority picker")
		}

	case !m.IsInModal() && key.Matches(msg, m.keyMap.ChangeType):
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			cmd = m.startChooseType(selected)
			m.logAction("tui opened type picker")
		}

	case !m.IsInModal() && key.Matches(msg, m.keyMap.ChangeAssignee):
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			cmd = m.startEditAssignee(selected)
			m.logAction("tui started editing assignee")
		}

	case !m.IsInModal() && key.Matches(msg, m.keyMap.AddComment):
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			cmd = m.startAddComment(selected)
		}

	case !m.IsInModal() && key.Matches(msg, m.keyMap.AddIssue):
		cmd = m.startCreateIssue()
		m.logAction("tui started creating issue")

	case !m.IsInModal() && key.Matches(msg, m.keyMap.DeleteIssue):
		fl := m.issueList
		if selected := fl.SelectedItem(); selected.ID != "" {
			cmd = m.startConfirmDelete(selected.ID, fl.Index())
			m.logAction("tui opened delete confirmation")
		}
	}

	return cmd
}
