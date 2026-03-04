package task

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/style"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type TaskDetails = models.TaskDetails

type TaskModel struct {
	TaskDetails
	keys          TaskHelpKeys
	width, height int
	userQuit      bool
}

type TaskHelpKeys struct {
	Quit  key.Binding
	Start key.Binding
}

var DefaultTaskKeys = TaskHelpKeys{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "Quit"),
	),
	Start: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Start"),
	),
}

func NewTaskModel(details TaskDetails) TaskModel {
	return TaskModel{
		TaskDetails: details,
		keys:        DefaultTaskKeys,
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
			m.userQuit = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Start):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m TaskModel) View() string {
	if m.width < 55 || m.height < 16 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			style.TextStyle.Render("Terminal too small."))
	}

	boxWidth := min(m.width-10, 120)

	detailsText := fmt.Sprintf("Time to complete: %s | Difficulty: %s", m.TimeToComplete, m.Difficulty)

	boxStyle := style.BorderStyle.
		Margin(1, 0).Padding(2, 4).Width(boxWidth)

	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(m.Title))
	b.WriteString("\n")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		style.TextStyle.Render(m.Description),
		"\n",
		style.TextStyle.Foreground(style.SecondaryColor).Render(detailsText),
	)

	b.WriteString(boxStyle.Render(content))
	b.WriteString("\n")

	helpText := "Press " + m.keys.Start.Help().Key + " to start • " + m.keys.Quit.Help().Key + " to quit"
	b.WriteString(style.HelpStyle.Render(helpText))

	final := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String())

	return final
}

func (m *TaskModel) SetSize(width, height int) {
	m.width, m.height = width, height
}

func (m TaskModel) GetUserQuit() bool {
	return m.userQuit
}
