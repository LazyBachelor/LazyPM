package dashboard

import (
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type DashboardKeyMap struct {
	components.CommonKeyMap
	SwitchWindow        key.Binding
	SwitchToKanbanBoard key.Binding
}

var defaultDashboardKeyMap = DashboardKeyMap{
	CommonKeyMap: components.DefaultCommonKeyMap(),
	SwitchWindow: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch window"),
	),
	SwitchToKanbanBoard: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "switch to kanban"),
	),
}

func (d *Model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, d.keyMap.Help):
		d.helpBar.ToggleHelp()
	case key.Matches(msg, d.keyMap.Quit):
		return tea.Quit
	case key.Matches(msg, d.keyMap.SwitchWindow):
		d.ToggleFocusedWindow()
	case !d.IsInModal() && key.Matches(msg, d.keyMap.SwitchToKanbanBoard):
		return func() tea.Msg { return msgs.SwitchToKanbanBoardMsg{} }
	case d.IsFocusedOnList() && key.Matches(msg, d.keyMap.SelectIssue):
		d.FocusDetail()
	case d.IsFocusedOnDetail() && (key.Matches(msg, d.keyMap.BackToList) || key.Matches(msg, d.keyMap.SelectIssue)):
		d.FocusList()
	case d.IsFocusedOnDetail() && key.Matches(msg, d.keyMap.ScrollUp):
		d.issueDetail.ScrollUp(1)
	case d.IsFocusedOnDetail() && key.Matches(msg, d.keyMap.ScrollDown):
		d.issueDetail.ScrollDown(1)
	case !d.IsInModal() && key.Matches(msg, d.keyMap.EditTitle):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startEditTitle(selected)
			cmd = d.titleInput.Focus()
		}
	case !d.IsInModal() && key.Matches(msg, d.keyMap.EditDescription):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startEditDescription(selected)
			cmd = d.descriptionInput.Focus()
		}
	case !d.IsInModal() && key.Matches(msg, d.keyMap.ChangeStatus):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChooseStatus(selected)
		}
	case !d.IsInModal() && key.Matches(msg, d.keyMap.ChangePriority):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChoosePriority(selected)
		}
	case !d.IsInModal() && key.Matches(msg, d.keyMap.ChangeType):
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
