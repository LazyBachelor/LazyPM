package dashboard

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
)

type DashboardKeyMap struct {
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

var defaultDashboardKeyMap = DashboardKeyMap{
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

func (d *Model) handleKeyPressMsg(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, d.keyMap.Help):
		d.helpBar.ToggleHelp()
		d.logAction("tui toggled help")
	case key.Matches(msg, d.keyMap.Quit):
		d.logAction("tui quit requested")
		return tea.Quit
	case !d.IsInModal() && key.Matches(msg, d.keyMap.SwitchToKanbanBoard):
		return func() tea.Msg { return msgs.SwitchToKanbanBoardMsg{} }
	case d.IsFocusedOnDetail() && key.Matches(msg, d.keyMap.ScrollUp):
		d.issueDetail.ScrollUp(1)
		d.logAction("tui scrolled issue detail up")
	case d.IsFocusedOnDetail() && key.Matches(msg, d.keyMap.ScrollDown):
		d.issueDetail.ScrollDown(1)
		d.logAction("tui scrolled issue detail down")
	case !d.IsInModal() && !d.addingComment && key.Matches(msg, d.keyMap.EditTitle):
		if selected := d.issueList.SelectedItem(); selected.ID != "" {
			d.startEditTitle(selected)
			cmd = d.titleInput.Focus()
			d.logAction("tui started editing issue title")
		}
	case !d.IsInModal() && !d.addingComment && key.Matches(msg, d.keyMap.EditDescription):
		if selected := d.issueList.SelectedItem(); selected.ID != "" {
			d.startEditDescription(selected)
			cmd = d.descriptionInput.Focus()
			d.logAction("tui started editing issue description")
		}
	case !d.IsInModal() && !d.addingComment && key.Matches(msg, d.keyMap.ChangeStatus):
		if selected := d.issueList.SelectedItem(); selected.ID != "" {
			d.startChooseStatus(selected)
			d.logAction("tui opened status picker")
		}
	case !d.IsInModal() && !d.addingComment && key.Matches(msg, d.keyMap.ChangePriority):
		if selected := d.issueList.SelectedItem(); selected.ID != "" {
			d.startChoosePriority(selected)
			d.logAction("tui opened priority picker")
		}
	case !d.IsInModal() && !d.addingComment && key.Matches(msg, d.keyMap.ChangeType):
		if selected := d.issueList.SelectedItem(); selected.ID != "" {
			d.startChooseType(selected)
			d.logAction("tui opened type picker")
		}
	case !d.IsInModal() && !d.addingComment && key.Matches(msg, d.keyMap.ChangeAssignee):
		if selected := d.issueList.SelectedItem(); selected.ID != "" {
			d.startEditAssignee(selected)
			cmd = d.assigneeInput.Focus()
			d.logAction("tui started editing assignee")
		}
	case !d.IsInModal() && !d.addingComment && key.Matches(msg, d.keyMap.AddComment):
		if selected := d.issueList.SelectedItem(); selected.ID != "" {
			d.startAddComment(selected)
			cmd = d.commentInput.Focus()
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && !d.editingAssignee && !d.addingComment && key.Matches(msg, d.keyMap.AddIssue):
		d.startCreateIssue()
		cmd = d.createTitleInput.Focus()
		d.logAction("tui started creating issue")
	case !d.IsInModal() && !d.addingComment && key.Matches(msg, d.keyMap.DeleteIssue):
		fl := d.issueList
		if selected := fl.SelectedItem(); selected.ID != "" {
			d.startConfirmDelete(selected.ID, fl.Index())
			d.logAction("tui opened delete confirmation")
		}
	}

	return cmd
}
