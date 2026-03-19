package task

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/style"
)

type TaskDetails = models.TaskDetails

type TaskModel struct {
	TaskDetails
	keys          TaskHelpKeys
	width, height int
	userQuit      bool
	aboutVisible  bool
}

type TaskHelpKeys struct {
	Quit  key.Binding
	Start key.Binding
	About key.Binding
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
	About: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Details about the interface"),
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
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.userQuit = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Start):
			return m, tea.Quit
		case key.Matches(msg, m.keys.About):
			m.aboutVisible = !m.aboutVisible
			return m, nil
		}
	}

	return m, nil
}

func (m TaskModel) View() tea.View {
	if m.width < 55 || m.height < 16 {
		content := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			style.TextStyle.Render("Terminal too small."))
		v := tea.NewView(content)
		v.AltScreen = true
		return v
	}

	boxWidth := min(m.width-10, 120)

	detailsText := fmt.Sprintf("Interface Type: %s | Time to complete: %s | Difficulty: %s", m.InterfaceType, m.TimeToComplete, m.Difficulty)

	boxStyle := style.BorderStyle.Padding(2, 4).Width(boxWidth)

	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(m.Title))
	b.WriteString("\n")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		style.TextStyle.Render(m.Description),
		"\n",
		style.TextStyle.Foreground(style.SecondaryColor).Render(detailsText),
		style.ErrorStyle.Render(m.interfaceHelpText()),
		style.ErrorStyle.Render(m.getQuitHelpText()),
	)

	if m.aboutVisible {
		b.WriteString(boxStyle.Render(m.InterfaceDescription))
	} else {
		b.WriteString(boxStyle.Render(content))

	}

	b.WriteString("\n")

	helpText := "Press " + m.keys.Start.Help().Key + " to start • " + m.keys.Quit.Help().Key + " to quit • " + m.keys.About.Help().Key + " " + m.keys.About.Help().Desc
	b.WriteString(style.HelpStyle.Render(helpText))

	final := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String())

	v := tea.NewView(final)
	v.AltScreen = true
	return v
}

func (m TaskModel) getQuitHelpText() string {
	switch m.InterfaceType {
	case models.InterfaceTypeWeb:
		return "Press esc/q in the terminal to skip the task"
	case models.InterfaceTypeTUI:
		return "Press 'q' to skip task during the survey"
	case models.InterfaceTypeCLI, models.InterfaceTypeREPL:
		return "Write 'exit' to skip task during the survey"
	default:
		return m.keys.Quit.Help().Key + " to quit"
	}
}

func (m TaskModel) interfaceHelpText() string {
	switch m.InterfaceType {
	case models.InterfaceTypeWeb:
		return "Press the button in the upper right corner to view task progress and completion criteria"
	case models.InterfaceTypeTUI:
		return "Press Shift+S to view task progress and completion criteria"
	case models.InterfaceTypeCLI, models.InterfaceTypeREPL:
		return "Write 'status' to view task progress and completion criteria"
	default:
		return m.keys.Quit.Help().Key + " to quit"
	}
}

func (m *TaskModel) SetSize(width, height int) {
	m.width, m.height = width, height
}

func (m TaskModel) GetUserQuit() bool {
	return m.userQuit
}
