package modal

import (
	"slices"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/style"
)

// TextInputResult is returned when a text input modal completes
type TextInputResult struct {
	Value string
}

// TextInputModal is a modal for single-line text input
// Suitable for: editing titles, editing assignees, creating issues
type TextInputModal struct {
	BaseModal
	input       textinput.Model
	label       string
	saveKeys    []string
	placeholder string
	charLimit   int
	width       int
	height      int
	issueID     string // Optional: for context
}

// TextInputConfig configures a text input modal
type TextInputConfig struct {
	ID           string
	Label        string
	Placeholder  string
	SaveKeys     []string
	CharLimit    int
	InitialValue string
	IssueID      string
	Width        int
	Height       int
}

// NewTextInputModal creates a new text input modal
func NewTextInputModal(cfg TextInputConfig) *TextInputModal {
	if cfg.SaveKeys == nil {
		cfg.SaveKeys = []string{"enter"}
	}
	if cfg.CharLimit == 0 {
		cfg.CharLimit = 256
	}

	ti := textinput.New()
	ti.Placeholder = cfg.Placeholder
	ti.CharLimit = cfg.CharLimit
	ti.SetValue(cfg.InitialValue)

	mod := &TextInputModal{
		BaseModal:   NewBaseModal(cfg.ID, TypeTextInput),
		input:       ti,
		label:       cfg.Label,
		saveKeys:    cfg.SaveKeys,
		placeholder: cfg.Placeholder,
		charLimit:   cfg.CharLimit,
		width:       cfg.Width,
		height:      cfg.Height,
		issueID:     cfg.IssueID,
	}

	if mod.width == 0 {
		mod.width = 60
	}
	if mod.height == 0 {
		mod.height = 20
	}

	return mod
}

// Activate prepares the modal for input
func (t *TextInputModal) Activate() tea.Cmd {
	t.BaseModal.activate()
	return t.input.Focus()
}

// Deactivate cleans up the modal
func (t *TextInputModal) Deactivate() {
	t.BaseModal.deactivate()
	t.input.Blur()
}

// SetValue updates the input value
func (t *TextInputModal) SetValue(value string) {
	t.input.SetValue(value)
}

// Value returns the current input value
func (t *TextInputModal) Value() string {
	return t.input.Value()
}

// IssueID returns the associated issue ID
func (t *TextInputModal) IssueID() string {
	return t.issueID
}

// Update handles input when the modal is active
func (t *TextInputModal) Update(msg tea.Msg) (tea.Cmd, bool) {
	if !t.IsActive() {
		return nil, false
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		s := msg.String()

		// Check save keys
		if slices.Contains(t.saveKeys, s) {
				value := t.input.Value()
				t.Deactivate()
				return func() tea.Msg {
					return ModalCompletedMsg{
						ModalID: t.ID(),
						Value:   TextInputResult{Value: value},
					}
				}, true
			}

		// Cancel on escape
		if s == "esc" {
			t.Deactivate()
			return func() tea.Msg {
				return ModalCancelledMsg{ModalID: t.ID()}
			}, true
		}
	}

	// Let the text input handle the message
	var cmd tea.Cmd
	t.input, cmd = t.input.Update(msg)
	// Always handle the message when modal is active to prevent keys leaking to list
	return cmd, true
}

// View renders the modal
func (t *TextInputModal) View() string {
	if t.width < 5 {
		return ""
	}

	boxWidth := min(60, t.width-4)
	t.input.SetWidth(boxWidth - 2)

	content := lipgloss.JoinVertical(lipgloss.Left,
		style.LabelStyle.Render(t.label),
		t.input.View(),
	)

	return style.ModalContainerStyle.
		Width(boxWidth).
		Render(content)
}

// SetSize updates the modal dimensions
func (t *TextInputModal) SetSize(width, height int) {
	t.BaseModal.SetSize(width, height)
	t.width = width
	t.height = height
}

// CursorEnd moves the cursor to the end of the input
func (t *TextInputModal) CursorEnd() {
	t.input.CursorEnd()
}

// Reset clears the input value
func (t *TextInputModal) Reset() {
	t.input.Reset()
}
