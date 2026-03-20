package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Manager provides a clean API for managing modals in views.
// Views should embed this struct to get modal management capabilities.
type Manager struct {
	stack  *ModalStack
	width  int
	height int
}

// NewManager creates a new modal manager
func NewManager() *Manager {
	return &Manager{
		stack: NewModalStack(),
	}
}

// SetSize updates the dimensions for modal rendering
func (m *Manager) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// ShowModal activates a pre-registered modal by ID
func (m *Manager) ShowModal(id string) tea.Cmd {
	for _, modal := range m.stack.modals {
		if modal.ID() == id {
			return modal.Activate()
		}
	}
	return nil
}

// RegisterModal adds a modal to the manager's registry
func (m *Manager) RegisterModal(modal Modal) {
	m.stack.Register(modal)
}

// PushModal adds a modal to the stack and activates it
func (m *Manager) PushModal(modal Modal) tea.Cmd {
	m.stack.Push(modal)
	return modal.Activate()
}

// PopModal removes the top modal from the stack
func (m *Manager) PopModal() Modal {
	return m.stack.Pop()
}

// CloseAll closes all active modals
func (m *Manager) CloseAll() {
	m.stack.Clear()
}

// IsModalActive returns true if any modal is currently active
func (m *Manager) IsModalActive() bool {
	return m.stack.HasActiveModal()
}

// ActiveModal returns the currently active modal
func (m *Manager) ActiveModal() Modal {
	return m.stack.ActiveModal()
}

// Update handles messages and routes them to the active modal
// Returns: (command, handled) - if handled is true, the view should stop processing this message
func (m *Manager) Update(msg tea.Msg) (tea.Cmd, bool) {
	return m.stack.Update(msg)
}

// View returns the rendered modal content (just the modal, not positioned)
func (m *Manager) View() string {
	active := m.stack.ActiveModal()
	if active == nil {
		return ""
	}
	active.SetSize(m.width, m.height)
	return active.View()
}

// GetTextInputModal retrieves a TextInputModal by ID (if registered)
func (m *Manager) GetTextInputModal(id string) *TextInputModal {
	for _, modal := range m.stack.modals {
		if modal.ID() == id && modal.Type() == TypeTextInput {
			if tim, ok := modal.(*TextInputModal); ok {
				return tim
			}
		}
	}
	return nil
}

// GetTextAreaModal retrieves a TextAreaModal by ID (if registered)
func (m *Manager) GetTextAreaModal(id string) *TextAreaModal {
	for _, modal := range m.stack.modals {
		if modal.ID() == id && modal.Type() == TypeTextArea {
			if tam, ok := modal.(*TextAreaModal); ok {
				return tam
			}
		}
	}
	return nil
}

// GetConfirmModal retrieves a ConfirmModal by ID (if registered)
func (m *Manager) GetConfirmModal(id string) *ConfirmModal {
	for _, modal := range m.stack.modals {
		if modal.ID() == id && modal.Type() == TypeConfirm {
			if cm, ok := modal.(*ConfirmModal); ok {
				return cm
			}
		}
	}
	return nil
}

// GetSelectModal retrieves a SelectModal by ID (if registered)
func (m *Manager) GetSelectModal(id string) *SelectModal {
	for _, modal := range m.stack.modals {
		if modal.ID() == id && modal.Type() == TypeSelect {
			if sm, ok := modal.(*SelectModal); ok {
				return sm
			}
		}
	}
	return nil
}

// RenderWithMainView renders the main view with an overlaid modal using Canvas
// This allows the modal to appear on top without clearing the background
func (m *Manager) RenderWithMainView(mainView string) string {
	if !m.IsModalActive() {
		return mainView
	}

	modalContent := m.View()
	if modalContent == "" {
		return mainView
	}

	// Calculate centered position for the modal
	modalWidth := lipgloss.Width(modalContent)
	modalHeight := lipgloss.Height(modalContent)

	// Center the modal
	x := max((m.width-modalWidth)/2, 0)
	y := max((m.height-modalHeight)/2, 0)

	// Create layers: main view as base, modal on top with Z-index
	mainLayer := lipgloss.NewLayer(mainView).X(0).Y(0).Z(0)
	modalLayer := lipgloss.NewLayer(modalContent).X(x).Y(y).Z(1)

	// Create compositor with layers and render
	compositor := lipgloss.NewCompositor(mainLayer, modalLayer)
	return compositor.Render()
}

// RegisterCommonModals registers the standard set of modals used across views.
// This helper reduces duplication between dashboard and kanban views.
func RegisterCommonModals(m *Manager) {

	// Exit Confirm Modal
	m.RegisterModal(NewConfirmModal(ConfirmConfig{
		ID:      ModalConfirmExit,
		Message: "Close the task and interface?",
		YesKeys: []string{"y", "Y"},
		NoKeys:  []string{"n", "N", "esc"},
	}))

	// Edit Title Modal
	m.RegisterModal(NewTextInputModal(TextInputConfig{
		ID:           ModalEditTitle,
		Label:        "Edit title (Enter to save, Esc to cancel):",
		Placeholder:  "Issue title...",
		SaveKeys:     []string{"enter"},
		CharLimit:    256,
		InitialValue: "",
	}))

	// Create Issue Modal
	m.RegisterModal(NewTextInputModal(TextInputConfig{
		ID:          ModalCreateIssue,
		Label:       "New issue (Enter to create, Esc to cancel):",
		Placeholder: "New issue title...",
		SaveKeys:    []string{"enter"},
		CharLimit:   256,
	}))

	// Edit Assignee Modal
	m.RegisterModal(NewTextInputModal(TextInputConfig{
		ID:          ModalEditAssignee,
		Label:       "Edit assignee (Enter to save, Esc to cancel):",
		Placeholder: "Assignee name...",
		SaveKeys:    []string{"enter"},
		CharLimit:   64,
	}))

	// Edit Description Modal
	m.RegisterModal(NewTextAreaModal(TextAreaConfig{
		ID:          ModalEditDescription,
		Label:       "Edit description (Ctrl+S to save, Esc to cancel):",
		Placeholder: "Issue description...",
		SaveKeys:    []string{"ctrl+s"},
		InputHeight: 10,
	}))

	// Add Comment Modal
	m.RegisterModal(NewTextAreaModal(TextAreaConfig{
		ID:          ModalAddComment,
		Label:       "Add comment (Ctrl+S to save, Esc to cancel):",
		Placeholder: "Write your comment...",
		SaveKeys:    []string{"ctrl+s"},
		InputHeight: 8,
	}))

	// Close Reason TextArea Modal
	m.RegisterModal(NewTextAreaModal(TextAreaConfig{
		ID:          ModalCloseReason,
		Label:       "Enter closing reason (Ctrl+S to save, Esc to cancel):",
		Placeholder: "Enter closing reason...",
		SaveKeys:    []string{"ctrl+s"},
		InputHeight: 4,
	}))

	// Delete Confirm Modal
	m.RegisterModal(NewConfirmModal(ConfirmConfig{
		ID:      ModalConfirmDelete,
		Message: "Delete issue?",
		YesKeys: []string{"y", "Y"},
		NoKeys:  []string{"n", "N", "esc"},
	}))

	// Status Select Modal
	m.RegisterModal(NewSelectModal(SelectConfig{
		ID:      ModalSelectStatus,
		Label:   "Change status:",
		Options: StatusOptions(),
	}))

	// Close Reason Select Modal
	m.RegisterModal(NewSelectModal(SelectConfig{
		ID:      ModalSelectCloseReason,
		Label:   "Choose closing reason:",
		Options: CloseReasonOptions(),
	}))

	// Priority Select Modal
	m.RegisterModal(NewSelectModal(SelectConfig{
		ID:      ModalSelectPriority,
		Label:   "Change priority:",
		Options: PriorityOptions(),
	}))

	// Type Select Modal
	m.RegisterModal(NewSelectModal(SelectConfig{
		ID:      ModalSelectType,
		Label:   "Change type:",
		Options: TypeOptions(),
	}))
}
