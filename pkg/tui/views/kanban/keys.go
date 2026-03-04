package kanban

import (
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KanbanKeyMap struct {
	components.CommonKeyMap
	SwitchToDashboard key.Binding
	MoveColumnLeft    key.Binding
	MoveColumnRight   key.Binding
	MoveIssueLeft     key.Binding
	MoveIssueRight    key.Binding
}

var defaultKanbanKeyMap = KanbanKeyMap{
	CommonKeyMap: components.DefaultCommonKeyMap(),
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

func (d *Model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, d.keyMap.Help):
		d.helpBar.ToggleHelp()
	case key.Matches(msg, d.keyMap.Quit):
		return tea.Quit
	case !d.IsInModal() && msg.String() == " ":
		return func() tea.Msg { return nil } // consume space
	case !d.IsInModal() && key.Matches(msg, d.keyMap.SwitchToDashboard):
		return func() tea.Msg { return msgs.SwitchToDashboardMsg{} }
	case !d.IsInModal() && key.Matches(msg, d.keyMap.MoveColumnLeft):
		if d.focusedColumn > 0 {
			d.focusedColumn--
			d.updateDetailFromSelection()
		}
	case !d.IsInModal() && key.Matches(msg, d.keyMap.MoveColumnRight):
		if d.focusedColumn < 2 {
			d.focusedColumn++
			d.updateDetailFromSelection()
		}
	case !d.IsInModal() && key.Matches(msg, d.keyMap.MoveIssueRight):
		cmd = d.moveIssue(+1)
	case !d.IsInModal() && key.Matches(msg, d.keyMap.MoveIssueLeft):
		cmd = d.moveIssue(-1)
	case d.IsFocusedOnList() && key.Matches(msg, d.keyMap.SelectIssue):
		d.FocusDetail()
	case d.IsFocusedOnDetail() && (key.Matches(msg, d.keyMap.BackToList) || key.Matches(msg, d.keyMap.SelectIssue)):
		d.FocusList()
	case d.IsFocusedOnDetail() && key.Matches(msg, d.keyMap.ScrollUp):
		d.issueDetail.ScrollUp(1)
	case d.IsFocusedOnDetail() && key.Matches(msg, d.keyMap.ScrollDown):
		d.issueDetail.ScrollDown(1)
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && key.Matches(msg, d.keyMap.EditTitle):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startEditTitle(selected)
			cmd = d.titleInput.Focus()
		}
	case !d.IsInModal() && key.Matches(msg, d.keyMap.EditDescription):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startEditDescription(selected)
			cmd = d.descriptionInput.Focus()
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && key.Matches(msg, d.keyMap.ChangeStatus):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChooseStatus(selected)
		}
	case !d.IsInModal() && key.Matches(msg, d.keyMap.ChangePriority):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChoosePriority(selected)
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && key.Matches(msg, d.keyMap.ChangeType):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChooseType(selected)
		}
	case !d.IsInModal() && key.Matches(msg, d.keyMap.AddIssue):
		d.startCreateIssue()
		cmd = d.createTitleInput.Focus()
	case !d.IsInModal() && key.Matches(msg, d.keyMap.DeleteIssue):
		fl := d.FocusedIssueList()
		if selected := fl.SelectedItem(); selected.ID != "" {
			d.startConfirmDelete(selected.ID, fl.Index())
		}
	}

	return cmd
}
