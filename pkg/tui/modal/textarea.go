package modal

import (
	"slices"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/style"
)

// TextAreaResult is returned when a text area modal completes
type TextAreaResult struct {
	Value string
}

// TextAreaModal is a modal for multi-line text input
// Suitable for: editing descriptions, adding comments, custom close reasons, etc.
type TextAreaModal struct {
	BaseModal
	input       textarea.Model
	label       string
	saveKeys    []string
	placeholder string
	width       int
	height      int
	inputHeight int
	issueID     string
}

// TextAreaConfig configures a text area modal
type TextAreaConfig struct {
	ID           string
	Label        string
	Placeholder  string
	SaveKeys     []string
	InitialValue string
	IssueID      string
	Width        int
	Height       int
	InputHeight  int // Height of the textarea itself
}

// NewTextAreaModal creates a new text area modal
func NewTextAreaModal(cfg TextAreaConfig) *TextAreaModal {
	if cfg.SaveKeys == nil {
		cfg.SaveKeys = []string{"ctrl+s"}
	}
	if cfg.InputHeight == 0 {
		cfg.InputHeight = 8
	}

	ta := textarea.New()
	ta.Placeholder = cfg.Placeholder
	ta.SetValue(cfg.InitialValue)

	mod := &TextAreaModal{
		BaseModal:   NewBaseModal(cfg.ID, TypeTextArea),
		input:       ta,
		label:       cfg.Label,
		saveKeys:    cfg.SaveKeys,
		placeholder: cfg.Placeholder,
		width:       cfg.Width,
		height:      cfg.Height,
		inputHeight: cfg.InputHeight,
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
func (t *TextAreaModal) Activate() tea.Cmd {
	t.BaseModal.activate()
	return t.input.Focus()
}

// Deactivate cleans up the modal
func (t *TextAreaModal) Deactivate() {
	t.BaseModal.deactivate()
	t.input.Blur()
}

// SetValue updates the textarea value
func (t *TextAreaModal) SetValue(value string) {
	t.input.SetValue(value)
}

// Value returns the current textarea value
func (t *TextAreaModal) Value() string {
	return t.input.Value()
}

// IssueID returns the associated issue ID
func (t *TextAreaModal) IssueID() string {
	return t.issueID
}

// Update handles input when the modal is active
func (t *TextAreaModal) Update(msg tea.Msg) (tea.Cmd, bool) {
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
						Value:   TextAreaResult{Value: value},
					}
				}, true
			}

		// Cancel on escape
		if s == "esc" {
			t.Deactivate()
			t.input.Blur()
			t.input.Reset()
			return func() tea.Msg {
				return ModalCancelledMsg{ModalID: t.ID()}
			}, true
		}
	}

	// Let the textarea handle the message
	var cmd tea.Cmd
	t.input, cmd = t.input.Update(msg)
	// Always handle the message when modal is active to prevent keys leaking to list
	return cmd, true
}

// View renders the modal
func (t *TextAreaModal) View() string {
	if t.width < 5 {
		return ""
	}

	boxWidth := max(min(60, t.width-4), 1)

	t.input.SetWidth(boxWidth - 2)
	t.input.SetHeight(t.inputHeight)

	content := lipgloss.JoinVertical(lipgloss.Left,
		style.LabelStyle.Render(t.label),
		t.input.View(),
	)

	return style.ModalContainerStyle.
		Width(boxWidth).
		Render(content)
}

// SetSize updates the modal dimensions
func (t *TextAreaModal) SetSize(width, height int) {
	t.BaseModal.SetSize(width, height)
	t.width = width
	t.height = height
}

// Reset clears the textarea
func (t *TextAreaModal) Reset() {
	t.input.Reset()
}
