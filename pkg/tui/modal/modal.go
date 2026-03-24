// Package modal provides a composable, interface-based modal system for TUI views.
// It enables separation of concerns between modal rendering, input handling, and focus management.
package modal

import (
	tea "charm.land/bubbletea/v2"
)

// Modal is the core interface that all modals must implement.
// It defines the contract for modal lifecycle, rendering, and input handling.
type Modal interface {
	ID() string
	Type() ModalType
	IsActive() bool
	Activate() tea.Cmd
	Deactivate()
	Update(msg tea.Msg) (tea.Cmd, bool)
	View() string
	SetSize(width, height int)
}

// ModalType categorizes different modal behaviors
type ModalType int

const (
	TypeTextInput ModalType = iota
	TypeConfirm
	TypeSelect
	TypeTextArea
	TypeCustom
)

// ModalResult carries the output from a completed modal
type ModalResult struct {
	ModalID string
	Value   interface{} // Type depends on the modal implementation
	Err     error
}

// ModalCompletedMsg is sent when a modal completes successfully
type ModalCompletedMsg struct {
	ModalID string
	Value   interface{}
}

// ModalCancelledMsg is sent when a modal is cancelled
type ModalCancelledMsg struct {
	ModalID string
}

// Modal IDs used across the application
const (
	ModalConfirmExit       = "confirm-exit"
	ModalEditTitle         = "edit-title"
	ModalCreateIssue       = "create-issue"
	ModalEditAssignee      = "edit-assignee"
	ModalEditDescription   = "edit-description"
	ModalAddComment        = "add-comment"
	ModalCloseReason       = "close-reason-other"
	ModalConfirmDelete     = "confirm-delete"
	ModalSelectStatus      = "select-status"
	ModalSelectCloseReason = "select-close-reason"
	ModalSelectPriority    = "select-priority"
	ModalSelectType        = "select-type"
	ModalSelectSprint       = "select-sprint"
	ModalManageDependencies   = "manage-dependencies"
	ModalSelectDependency     = "select-dependency"
	ModalSelectRemoveDependency = "select-remove-dependency"
)

// ModalStack manages a stack of active modals with priority handling
type ModalStack struct {
	modals []Modal
}

// NewModalStack creates an empty modal stack
func NewModalStack() *ModalStack {
	return &ModalStack{
		modals: make([]Modal, 0),
	}
}

// Push adds a modal to the top of the stack
func (s *ModalStack) Push(m Modal) {
	s.modals = append(s.modals, m)
}

// Pop removes and returns the top modal
func (s *ModalStack) Pop() Modal {
	if len(s.modals) == 0 {
		return nil
	}
	m := s.modals[len(s.modals)-1]
	s.modals = s.modals[:len(s.modals)-1]
	return m
}

// Peek returns the top modal without removing it
func (s *ModalStack) Peek() Modal {
	if len(s.modals) == 0 {
		return nil
	}
	return s.modals[len(s.modals)-1]
}

// ActiveModal returns the currently active modal (top of stack if active)
func (s *ModalStack) ActiveModal() Modal {
	for i := len(s.modals) - 1; i >= 0; i-- {
		if s.modals[i].IsActive() {
			return s.modals[i]
		}
	}
	return nil
}

// HasActiveModal returns true if any modal in the stack is active
func (s *ModalStack) HasActiveModal() bool {
	return s.ActiveModal() != nil
}

// Clear deactivates and removes all modals
func (s *ModalStack) Clear() {
	for _, m := range s.modals {
		m.Deactivate()
	}
	s.modals = s.modals[:0]
}

// Update routes messages to the active modal
// Returns handled=false for global messages like WindowSizeMsg so they pass through to the view
func (s *ModalStack) Update(msg tea.Msg) (tea.Cmd, bool) {
	active := s.ActiveModal()
	if active == nil {
		return nil, false
	}

	// Allow global quit key even when modal is active
	if key, ok := msg.(tea.KeyPressMsg); ok {
		if key.String() == "ctrl+c" {
			return nil, false
		}
	}

	// WindowSizeMsg should pass through to the view even when modal is active
	// This allows both the modal and the underlying view to resize
	if _, ok := msg.(tea.WindowSizeMsg); ok {
		cmd, _ := active.Update(msg)
		return cmd, false
	}

	return active.Update(msg)
}

// View renders the active modal with proper framing
func (s *ModalStack) View(width, height int) string {
	active := s.ActiveModal()
	if active == nil {
		return ""
	}
	active.SetSize(width, height)
	return active.View()
}

// Len returns the number of modals in the stack
func (s *ModalStack) Len() int {
	return len(s.modals)
}

// Register allows pre-registering modals for use
func (s *ModalStack) Register(modals ...Modal) {
	s.modals = append(s.modals, modals...)
}
