package kanban

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
)

type KeyMap struct {
	components.CommonKeyMap
	SwitchToDashboard key.Binding
	MoveColumnLeft    key.Binding
	MoveColumnRight   key.Binding
	MoveIssueLeft     key.Binding
	MoveIssueRight    key.Binding
	SubmitValidation  key.Binding
}

var defaultKanbanKeyMap = KeyMap{
	CommonKeyMap: components.DefaultCommonKeyMap(),
	SubmitValidation: key.NewBinding(
		key.WithKeys("S"),
		key.WithHelp("S", "submit validation"),
	),
	SwitchToDashboard: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "dashboard 1"),
	),
	MoveColumnLeft: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "prev column"),
	),
	MoveColumnRight: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "next column"),
	),
	MoveIssueLeft: key.NewBinding(
		key.WithKeys("left", "["),
		key.WithHelp("←/[", "move issue left"),
	),
	MoveIssueRight: key.NewBinding(
		key.WithKeys("right", "]"),
		key.WithHelp("→/]", "move issue right"),
	),
}

func (m *Model) handleKeyPressMsg(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keyMap.SubmitValidation):
		if m.submitChan != nil {
			select {
			case m.submitChan <- struct{}{}:
			default:
			}
		}
	case key.Matches(msg, m.keyMap.Help):
		m.helpBar.ToggleHelp()
	case key.Matches(msg, m.keyMap.Quit):
		return tea.Quit
	case !m.IsInModal() && msg.String() == " ":
		return func() tea.Msg { return nil }
	case !m.IsInModal() && key.Matches(msg, m.keyMap.SwitchToDashboard):
		return func() tea.Msg { return msgs.SwitchToDashboardMsg{} }
	case !m.IsInModal() && key.Matches(msg, m.keyMap.MoveColumnLeft):
		m.focusManager.PreviousColumn()
		m.updateDetailFromSelection()
	case !m.IsInModal() && key.Matches(msg, m.keyMap.MoveColumnRight):
		m.focusManager.NextColumn()
		m.updateDetailFromSelection()
	case !m.IsInModal() && key.Matches(msg, m.keyMap.MoveIssueRight):
		cmd = m.moveIssue(+1)
	case !m.IsInModal() && key.Matches(msg, m.keyMap.MoveIssueLeft):
		cmd = m.moveIssue(-1)
	case m.IsFocusedOnDetail() && key.Matches(msg, m.keyMap.ScrollUp):
		m.issueDetail.ScrollUp(1)
	case m.IsFocusedOnDetail() && key.Matches(msg, m.keyMap.ScrollDown):
		m.issueDetail.ScrollDown(1)
	case !m.IsInModal() && key.Matches(msg, m.keyMap.EditTitle):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startEditTitle(selected)
		}
	case !m.IsInModal() && key.Matches(msg, m.keyMap.EditDescription):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startEditDescription(selected)
		}
	case !m.IsInModal() && key.Matches(msg, m.keyMap.ChangeStatus):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startChooseStatus(selected)
		}
	case !m.IsInModal() && key.Matches(msg, m.keyMap.ChangePriority):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startChoosePriority(selected)
		}
	case !m.IsInModal() && key.Matches(msg, m.keyMap.ChangeType):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startChooseType(selected)
		}
	case !m.IsInModal() && key.Matches(msg, m.keyMap.ChangeAssignee):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startEditAssignee(selected)
		}
	case !m.IsInModal() && key.Matches(msg, m.keyMap.AddIssue):
		cmd = m.startCreateIssue()
	case !m.IsInModal() && key.Matches(msg, m.keyMap.DeleteIssue):
		fl := m.FocusedIssueList()
		if selected := fl.SelectedItem(); selected.ID != "" {
			cmd = m.startConfirmDelete(selected.ID, fl.Index())
		}
	}

	return cmd
}
