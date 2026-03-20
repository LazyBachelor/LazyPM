package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
)

// SelectResult is returned when a select modal completes
type SelectResult struct {
	SelectedKey   string
	SelectedValue string
}

// SelectOption represents a selectable option
type SelectOption struct {
	Key   string // The key to press
	Label string // Display label
	Value string // The actual value to return
}

// SelectModal is a modal for selecting from a list of options
// Suitable for: status selection, priority selection, type selection, etc.
type SelectModal struct {
	BaseModal
	label     string
	options   []SelectOption
	helpText  string
	cancelKey string
	issueID   string
	width     int
	height    int
}

// SelectConfig configures a select modal
type SelectConfig struct {
	ID        string
	Label     string
	Options   []SelectOption
	CancelKey string
	IssueID   string
	Width     int
	Height    int
}

// NewSelectModal creates a new select modal
func NewSelectModal(cfg SelectConfig) *SelectModal {
	if cfg.CancelKey == "" {
		cfg.CancelKey = "esc"
	}

	mod := &SelectModal{
		BaseModal: NewBaseModal(cfg.ID, TypeSelect),
		label:     cfg.Label,
		options:   cfg.Options,
		cancelKey: cfg.CancelKey,
		issueID:   cfg.IssueID,
		width:     cfg.Width,
		height:    cfg.Height,
	}

	// Build help text from options
	var parts []string
	for _, opt := range cfg.Options {
		parts = append(parts, opt.Key+" = "+opt.Label)
	}
	mod.helpText = lipgloss.NewStyle().Foreground(styles.FaintText).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, parts...))

	if mod.width == 0 {
		mod.width = 70
	}
	if mod.height == 0 {
		mod.height = 20
	}

	return mod
}

// Activate prepares the modal
func (s *SelectModal) Activate() tea.Cmd {
	s.BaseModal.activate()
	return nil
}

// Deactivate cleans up the modal
func (s *SelectModal) Deactivate() {
	s.BaseModal.deactivate()
}

// IssueID returns the associated issue ID
func (s *SelectModal) IssueID() string {
	return s.issueID
}

// Options returns the available options
func (s *SelectModal) Options() []SelectOption {
	return s.options
}

// Update handles input when the modal is active
func (s *SelectModal) Update(msg tea.Msg) (tea.Cmd, bool) {
	if !s.IsActive() {
		return nil, false
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()

		// Check cancel key
		if key == s.cancelKey {
			s.Deactivate()
			return func() tea.Msg {
				return ModalCancelledMsg{ModalID: s.ID()}
			}, true
		}

		// Check option keys
		for _, opt := range s.options {
			if key == opt.Key {
				s.Deactivate()
				return func() tea.Msg {
					return ModalCompletedMsg{
						ModalID: s.ID(),
						Value:   SelectResult{SelectedKey: opt.Key, SelectedValue: opt.Value},
					}
				}, true
			}
		}
	}

	// Consume all keys when modal is active to prevent leakage to underlying components
	return nil, true
}

// View renders the modal
func (s *SelectModal) View() string {
	if s.width < 5 {
		return ""
	}

	boxWidth := min(70, s.width-4)
	if boxWidth < 1 {
		boxWidth = 1
	}

	cancelText := lipgloss.NewStyle().
		Foreground(styles.FaintText).
		Render(s.cancelKey + " = cancel")

	content := lipgloss.JoinVertical(lipgloss.Left,
		styles.LabelStyle.Render(s.label),
		s.helpText,
		cancelText,
	)

	return styles.ModalContainerStyle.
		Width(boxWidth).
		Render(content)
}

// SetSize updates the modal dimensions
func (s *SelectModal) SetSize(width, height int) {
	s.BaseModal.SetSize(width, height)
	s.width = width
	s.height = height
}

// Predefined option sets for common use cases

// StatusOptions returns options for status selection
func StatusOptions() []SelectOption {
	return []SelectOption{
		{Key: "o", Label: "open", Value: "open"},
		{Key: "i", Label: "in_progress", Value: "in_progress"},
		{Key: "b", Label: "blocked", Value: "blocked"},
		{Key: "r", Label: "ready_to_sprint", Value: "ready_to_sprint"},
		{Key: "c", Label: "closing", Value: "closing"},
	}
}

// PriorityOptions returns options for priority selection
func PriorityOptions() []SelectOption {
	return []SelectOption{
		{Key: "0", Label: "irrelevant", Value: "0"},
		{Key: "1", Label: "low", Value: "1"},
		{Key: "2", Label: "normal", Value: "2"},
		{Key: "3", Label: "high", Value: "3"},
		{Key: "4", Label: "critical", Value: "4"},
	}
}

// TypeOptions returns options for issue type selection
func TypeOptions() []SelectOption {
	return []SelectOption{
		{Key: "b", Label: "bug", Value: "bug"},
		{Key: "f", Label: "feature", Value: "feature"},
		{Key: "t", Label: "task", Value: "task"},
		{Key: "e", Label: "epic", Value: "epic"},
		{Key: "c", Label: "chore", Value: "chore"},
	}
}

// CloseReasonOptions returns options for close reason selection
func CloseReasonOptions() []SelectOption {
	return []SelectOption{
		{Key: "d", Label: "Done", Value: "Done"},
		{Key: "u", Label: "Duplicate issue", Value: "Duplicate issue"},
		{Key: "w", Label: "Won't fix", Value: "Won't fix"},
		{Key: "o", Label: "Obsolete", Value: "Obsolete"},
		{Key: "h", Label: "Other", Value: "other"},
	}
}
