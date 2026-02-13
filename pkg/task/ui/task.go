package taskui

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func NewTaskModel(details TaskDetails) TaskModel {
	return TaskModel{
		TaskDetails: details,
		keys:        DefaultTaskKeys,
		help:        help.New(),
	}
}

func (m TaskModel) Init() tea.Cmd {
	return nil
}

func (m TaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Continue):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m TaskModel) View() string {
	padding := 3

	header := lipgloss.NewStyle().
		PaddingTop(padding).Width(m.width).Align(lipgloss.Center).
		Bold(true).Render(m.Title)

	headerHeight := lipgloss.Height(header)

	helpView := lipgloss.NewStyle().
		PaddingBottom(padding).
		Width(m.width).Align(lipgloss.Center).
		Render(m.help.View(m.keys))

	helpHeigh := lipgloss.Height(helpView)

	detailsText := fmt.Sprintf("Time to complete: %s | Difficulty: %s", m.TimeToComplete, m.Difficulty)
	details := lipgloss.NewStyle().Align(lipgloss.Center).
		Width(m.width).PaddingBottom(1).Render(detailsText)

	detailsHeight := lipgloss.Height(details)

	content := lipgloss.NewStyle().
		Width(m.width).Height(m.height-headerHeight-helpHeigh-detailsHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(m.Description)

	return lipgloss.JoinVertical(lipgloss.Top, header, content, details, helpView)

}

func (m *TaskModel) SetSize(width, height int) {
	m.width, m.height = width, height
}
