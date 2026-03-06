package dashboard

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type DashboardKeyMap struct {
	Help             key.Binding
	Quit             key.Binding
	SelectIssue      key.Binding
	BackToList       key.Binding
	ScrollUp         key.Binding
	ScrollDown       key.Binding
	SwitchWindow     key.Binding
	EditTitle        key.Binding
	EditDescription  key.Binding
	ChangeStatus     key.Binding
	ChangePriority   key.Binding
	ChangeType       key.Binding
	AddComment       key.Binding
	AddIssue         key.Binding
	DeleteIssue      key.Binding
	SubmitValidation key.Binding
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
	SubmitValidation: key.NewBinding(
		key.WithKeys("S"),
		key.WithHelp("S", "submit validation"),
	),
}

func (d *Model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, d.keyMap.Help):
		d.helpBar.ToggleHelp()
		d.logAction("tui toggled help")
	case key.Matches(msg, d.keyMap.Quit):
		d.logAction("tui quit requested")
		return tea.Quit
	case key.Matches(msg, d.keyMap.SwitchWindow):
		d.ToggleFocusedWindow()
		d.logAction("tui switched issue window")
	case d.IsFocusedOnList() && key.Matches(msg, d.keyMap.SelectIssue):
		d.FocusDetail()
		d.logAction("tui opened issue detail")
	case d.IsFocusedOnDetail() && key.Matches(msg, d.keyMap.BackToList):
		d.FocusList()
		d.logAction("tui returned to issue list")
	case d.IsFocusedOnDetail() && key.Matches(msg, d.keyMap.ScrollUp):
		d.issueDetail.ScrollUp(1)
		d.logAction("tui scrolled issue detail up")
	case d.IsFocusedOnDetail() && key.Matches(msg, d.keyMap.ScrollDown):
		d.issueDetail.ScrollDown(1)
		d.logAction("tui scrolled issue detail down")
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && !d.addingComment && key.Matches(msg, d.keyMap.EditTitle):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startEditTitle(selected)
			cmd = d.titleInput.Focus()
			d.logAction("tui started editing issue title")
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && !d.addingComment && key.Matches(msg, d.keyMap.EditDescription):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startEditDescription(selected)
			cmd = d.descriptionInput.Focus()
			d.logAction("tui started editing issue description")
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && !d.addingComment && key.Matches(msg, d.keyMap.ChangeStatus):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChooseStatus(selected)
			d.logAction("tui opened status picker")
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && !d.addingComment && key.Matches(msg, d.keyMap.ChangePriority):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChoosePriority(selected)
			d.logAction("tui opened priority picker")
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && !d.addingComment && key.Matches(msg, d.keyMap.ChangeType):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startChooseType(selected)
			d.logAction("tui opened type picker")
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && !d.addingComment && key.Matches(msg, d.keyMap.AddComment):
		if selected := d.FocusedIssueList().SelectedItem(); selected.ID != "" {
			d.startAddComment(selected)
			cmd = d.commentInput.Focus()
		}
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && !d.addingComment && key.Matches(msg, d.keyMap.AddIssue):
		d.startCreateIssue()
		cmd = d.createTitleInput.Focus()
		d.logAction("tui started creating issue")
	case !d.editingTitle && !d.creatingIssue && !d.editingDescription && !d.choosingStatus && !d.choosingPriority && !d.confirmingDelete && !d.choosingType && !d.addingComment && key.Matches(msg, d.keyMap.DeleteIssue):
		fl := d.FocusedIssueList()
		if selected := fl.SelectedItem(); selected.ID != "" {
			d.startConfirmDelete(selected.ID, fl.Index())
			d.logAction("tui opened delete confirmation")
		}
	case key.Matches(msg, d.keyMap.SubmitValidation):
		if d.submitChan != nil {
			select {
			case d.submitChan <- struct{}{}:
				d.logAction("tui submitted validation")
			default:
			}
		}
	}

	return cmd
}
