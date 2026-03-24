package components

import "charm.land/bubbles/v2/key"

// CommonKeyMap holds the key bindings shared between the issues dashboard
// and the kanban board.
type CommonKeyMap struct {
	Help            key.Binding
	Quit            key.Binding
	ScrollUp        key.Binding
	ScrollDown      key.Binding
	EditTitle       key.Binding
	EditDescription key.Binding
	ChangeStatus    key.Binding
	ChangePriority  key.Binding
	ChangeType      key.Binding
	ChangeAssignee    key.Binding
	ManageDependencies key.Binding
	AddIssue          key.Binding
	DeleteIssue       key.Binding
}

// DefaultCommonKeyMap returns the shared default bindings used by both views.
func DefaultCommonKeyMap() CommonKeyMap {
	return CommonKeyMap{
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		ScrollUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
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
		ChangeAssignee: key.NewBinding(
			key.WithKeys("A"),
			key.WithHelp("A", "change assignee"),
		),
		ManageDependencies: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "edit dependencies"),
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
}
