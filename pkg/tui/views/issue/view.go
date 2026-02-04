package issue

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Viewport      viewport.Model
	ID            string
	Title         string
	Description   string
	Status        string
	IssueType     string
	Width, Height int
	ready         bool
}

func NewIssueView(issue Model) Model {
	return Model{
		ID:          issue.ID,
		Title:       issue.Title,
		Description: issue.Description,
		Status:      issue.Status,
		IssueType:   issue.IssueType,
		Width:       issue.Width,
		Height:      issue.Height,
		Viewport:    viewport.New(issue.Width, issue.Height),
		ready:       false,
	}
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		m.Viewport.Width = m.Width
		m.Viewport.Height = m.Height
		m.ready = true
	}

	if !m.ready {
		return m, nil
	}

	content := fmt.Sprintf("ID: %s\nTitle: %s\nDescription: %s\nStatus: %s\nType: %s",
		m.ID, m.Title, m.Description, m.Status, m.IssueType)

	m.Viewport.SetContent(content)
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("%s", m.Viewport.View())
}
