package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
)

// ConfirmResult is returned when a confirm modal completes
type ConfirmResult struct {
	Confirmed bool
}

// ConfirmModal is a modal for yes/no confirmations
// Suitable for: delete confirmations, discard changes, etc.
type ConfirmModal struct {
	BaseModal
	message string
	yesKeys []string
	noKeys  []string
	issueID string
	width   int
	height  int
}

// ConfirmConfig configures a confirmation modal
type ConfirmConfig struct {
	ID      string
	Message string
	YesKeys []string
	NoKeys  []string
	IssueID string
	Width   int
	Height  int
}

// NewConfirmModal creates a new confirmation modal
func NewConfirmModal(cfg ConfirmConfig) *ConfirmModal {
	if cfg.YesKeys == nil {
		cfg.YesKeys = []string{"y", "Y"}
	}
	if cfg.NoKeys == nil {
		cfg.NoKeys = []string{"n", "N", "esc"}
	}

	mod := &ConfirmModal{
		BaseModal: NewBaseModal(cfg.ID, TypeConfirm),
		message:   cfg.Message,
		yesKeys:   cfg.YesKeys,
		noKeys:    cfg.NoKeys,
		issueID:   cfg.IssueID,
		width:     cfg.Width,
		height:    cfg.Height,
	}

	if mod.width == 0 {
		mod.width = 50
	}
	if mod.height == 0 {
		mod.height = 20
	}

	return mod
}

// Activate prepares the modal
func (c *ConfirmModal) Activate() tea.Cmd {
	c.BaseModal.activate()
	return nil
}

// Deactivate cleans up the modal
func (c *ConfirmModal) Deactivate() {
	c.BaseModal.deactivate()
}

// IssueID returns the associated issue ID
func (c *ConfirmModal) IssueID() string {
	return c.issueID
}

// Update handles input when the modal is active
func (c *ConfirmModal) Update(msg tea.Msg) (tea.Cmd, bool) {
	if !c.IsActive() {
		return nil, false
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		s := msg.String()

		// Check yes keys
		for _, key := range c.yesKeys {
			if s == key {
				c.Deactivate()
				return func() tea.Msg {
					return ModalCompletedMsg{
						ModalID: c.ID(),
						Value:   ConfirmResult{Confirmed: true},
					}
				}, true
			}
		}

		// Check no/cancel keys
		for _, key := range c.noKeys {
			if s == key {
				c.Deactivate()
				return func() tea.Msg {
					return ModalCancelledMsg{ModalID: c.ID()}
				}, true
			}
		}
	}

	// Consume all keys when modal is active to prevent leakage to underlying components
	return nil, true
}

// View renders the modal
func (c *ConfirmModal) View() string {
	if c.width < 5 {
		return ""
	}

	boxWidth := min(50, c.width-4)
	if boxWidth < 1 {
		boxWidth = 1
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		styles.LabelStyle.Render(c.message),
		lipgloss.NewStyle().Foreground(styles.FaintText).
			Render("Press y to confirm, n or Esc to cancel"),
	)

	return styles.ModalContainerStyle.
		Width(boxWidth).
		Render(content)
}

// SetSize updates the modal dimensions
func (c *ConfirmModal) SetSize(width, height int) {
	c.BaseModal.SetSize(width, height)
	c.width = width
	c.height = height
}
