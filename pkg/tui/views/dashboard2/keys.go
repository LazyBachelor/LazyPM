package dashboard2

import (
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type DashboardKeyMap struct {
	Help              key.Binding
	Quit              key.Binding
	SelectIssue       key.Binding
	BackToList        key.Binding
	ScrollUp          key.Binding
	ScrollDown        key.Binding
	SwitchWindow      key.Binding
	SwitchToDashboard key.Binding
	EditTitle         key.Binding
	EditDescription   key.Binding
	ChangeStatus      key.Binding
	ChangePriority    key.Binding
	ChangeType        key.Binding
	AddIssue          key.Binding
	DeleteIssue       key.Binding
}

var defaultDashboardKeyMap = DashboardKeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
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
		key.WithHelp("↑/k", "up"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	SwitchWindow: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch window"),
	),
	SwitchToDashboard: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "dashboard 1"),
	),
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
	AddIssue: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add issue"),
	),
	DeleteIssue: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete issue"),
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
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && key.Matches(msg, d.keyMap.SwitchToDashboard):
		return func() tea.Msg { return msgs.SwitchToDashboardMsg{} }
	case d.IsFocusedOnList() && key.Matches(msg, d.keyMap.SelectIssue):
		d.FocusDetail()
	case d.IsFocusedOnDetail() && key.Matches(msg, d.keyMap.BackToList):
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
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && key.Matches(msg, d.keyMap.EditDescription):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startEditDescription(selected)
			cmd = d.descriptionInput.Focus()
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && key.Matches(msg, d.keyMap.ChangeStatus):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChooseStatus(selected)
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && key.Matches(msg, d.keyMap.ChangePriority):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChoosePriority(selected)
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && key.Matches(msg, d.keyMap.ChangeType):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChooseType(selected)
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && key.Matches(msg, d.keyMap.AddIssue):
		d.startCreateIssue()
		cmd = d.createTitleInput.Focus()
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && key.Matches(msg, d.keyMap.DeleteIssue):
		fl := d.FocusedIssueList()
		if selected := fl.SelectedItem(); selected.ID != "" {
			d.startConfirmDelete(selected.ID, fl.Index())
		}
	}

	return cmd
}
