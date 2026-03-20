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
	AddComment        key.Binding
}

var defaultKanbanKeyMap = KeyMap{
	CommonKeyMap: components.DefaultCommonKeyMap(),
	SwitchToDashboard: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "dashboard"),
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
	AddComment: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "add comment"),
	),
}

func (m *Model) handleKeyPressMsg(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case m.notInModalMsgWithKey(msg, m.keyMap.Help):
		m.helpBar.ToggleHelp()

	case m.notInModalMsgWithKey(msg, m.keyMap.Quit):
		return tea.Quit

	case m.notInModalMsgWithKey(msg, m.keyMap.SwitchToDashboard):
		return func() tea.Msg { return msgs.SwitchToDashboardMsg{} }

	case m.notInModalMsgWithKey(msg, m.keyMap.MoveColumnLeft):
		m.focusManager.PreviousColumn()
		m.updateDetailFromSelection()

	case m.notInModalMsgWithKey(msg, m.keyMap.MoveColumnRight):
		m.focusManager.NextColumn()
		m.updateDetailFromSelection()

	case m.notInModalMsgWithKey(msg, m.keyMap.MoveIssueRight):
		cmd = m.moveIssue(+1)

	case m.notInModalMsgWithKey(msg, m.keyMap.MoveIssueLeft):
		cmd = m.moveIssue(-1)

	case m.IsFocusedOnDetail() && key.Matches(msg, m.keyMap.ScrollUp):
		m.issueDetail.ScrollUp(1)

	case m.notInModalMsgWithKey(msg, m.keyMap.ScrollDown):
		m.issueDetail.ScrollDown(1)

	case m.notInModalMsgWithKey(msg, m.keyMap.EditTitle):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startEditTitle(selected)
		}
	case m.notInModalMsgWithKey(msg, m.keyMap.EditDescription):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startEditDescription(selected)
		}
	case m.notInModalMsgWithKey(msg, m.keyMap.ChangeStatus):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startChooseStatus(selected)
		}
	case m.notInModalMsgWithKey(msg, m.keyMap.ChangePriority):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startChoosePriority(selected)
		}
	case m.notInModalMsgWithKey(msg, m.keyMap.ChangeType):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startChooseType(selected)
		}
	case m.notInModalMsgWithKey(msg, m.keyMap.ChangeAssignee):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startEditAssignee(selected)
		}
	case m.notInModalMsgWithKey(msg, m.keyMap.AddIssue):
		cmd = m.startCreateIssue()

	case m.notInModalMsgWithKey(msg, m.keyMap.AddComment):
		if selected := m.FocusedIssueList().SelectedItem(); selected.ID != "" {
			cmd = m.startAddComment(selected)
		}

	case m.notInModalMsgWithKey(msg, m.keyMap.DeleteIssue):
		fl := m.FocusedIssueList()
		if selected := fl.SelectedItem(); selected.ID != "" {
			cmd = m.startConfirmDelete(selected.ID, fl.Index())
		}
	}

	return cmd
}

func (m *Model) notInModalMsgWithKey(msg tea.KeyPressMsg, keyBinding key.Binding) bool {
	return !m.IsInModal() && key.Matches(msg, keyBinding)
}
