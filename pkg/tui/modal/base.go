package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/style"
)

// BaseModal provides common functionality for all modals.
// Embed this struct to get default implementations of the Modal interface.
type BaseModal struct {
	id      string
	modType ModalType
	active  bool
	width   int
	height  int
}

// NewBaseModal creates a new base modal with the given properties
func NewBaseModal(id string, modType ModalType) BaseModal {
	return BaseModal{
		id:      id,
		modType: modType,
		width:   80,
		height:  20,
	}
}

// ID returns the modal's unique identifier
func (b *BaseModal) ID() string {
	return b.id
}

// Type returns the modal type
func (b *BaseModal) Type() ModalType {
	return b.modType
}

// IsActive returns true if the modal is currently active
func (b *BaseModal) IsActive() bool {
	return b.active
}

// SetSize updates the modal dimensions
func (b *BaseModal) SetSize(width, height int) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	b.width = width
	b.height = height
}

// activate marks the modal as active
func (b *BaseModal) activate() tea.Cmd {
	b.active = true
	return nil
}

// deactivate marks the modal as inactive
func (b *BaseModal) deactivate() {
	b.active = false
}

// Width returns the modal width
func (b *BaseModal) Width() int {
	return b.width
}

// Height returns the modal height
func (b *BaseModal) Height() int {
	return b.height
}

// ModalFrame renders content within a standard modal frame without full-screen placement
func ModalFrame(content string, width int) string {
	if width < 5 {
		return ""
	}

	boxWidth := max(min(80, width-4), 1)

	return style.ModalContainerStyle.
		Width(boxWidth).
		Render(content)
}

// ModalWithLabel renders a modal with a label and content
func ModalWithLabel(label, content string, width int) string {
	if width < 5 {
		return ""
	}

	fullContent := lipgloss.JoinVertical(lipgloss.Left,
		style.LabelStyle.Render(label),
		content,
	)

	return ModalFrame(fullContent, width)
}
