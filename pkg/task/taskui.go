package task

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
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
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Start: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "start"),
	),
	About: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "interface info"),
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
		content := lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			style.TextStyle.Render("Terminal too small."),
		)
		v := tea.NewView(content)
		v.AltScreen = true
		return v
	}

	boxWidth := min(m.width-8, 120)

	cardStyle := style.BorderStyle.Padding(1, 2)
	bodyWidth := boxWidth - cardStyle.GetHorizontalFrameSize()
	cardStyle = cardStyle.Width(bodyWidth)

	textWidth := min(bodyWidth, 64)

	titleWrap := lipgloss.NewStyle().
		Width(boxWidth).
		Align(lipgloss.Center)

	helpWrap := lipgloss.NewStyle().
		Width(boxWidth).
		Align(lipgloss.Center)

	descStyle := style.TextStyle.
		Width(textWidth).
		Align(lipgloss.Left)

	detailsStyle := style.TextStyle.
		Foreground(style.SecondaryText).
		Width(textWidth).
		Align(lipgloss.Center)

	infoStyle := style.ErrorStyle.
		Width(textWidth).
		Align(lipgloss.Center)

	centerInCard := func(s string) string {
		return lipgloss.PlaceHorizontal(bodyWidth, lipgloss.Center, s)
	}

	detailsText := fmt.Sprintf(
		"Interface Type: %s | Time to complete: %s | Difficulty: %s",
		m.InterfaceType,
		m.TimeToComplete,
		m.Difficulty,
	)

	var body string
	if m.aboutVisible {
		about := style.TextStyle.
			Width(bodyWidth - 2).
			Align(lipgloss.Left).
			Render(m.InterfaceDescription)

		body = centerInCard(about)
	} else {
		body = lipgloss.JoinVertical(
			lipgloss.Left,
			centerInCard(descStyle.Render(m.Description)),
			"",
			centerInCard(detailsStyle.Render(detailsText)),
			centerInCard(infoStyle.Render(m.interfaceHelpText())),
			centerInCard(infoStyle.Render(m.getQuitHelpText())),
		)
	}

	card := cardStyle.Render(body)

	helpText := "Press " + m.keys.Start.Help().Key + " to " + m.keys.Start.Help().Desc +
		"  • " + m.keys.Quit.Help().Key + " to " + m.keys.Quit.Help().Desc +
		" • " + m.keys.About.Help().Key + " for " + m.keys.About.Help().Desc

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		titleWrap.Render(style.TitleStyle.Render(m.Title)),
		"",
		card,
		"",
		helpWrap.Render(style.HelpStyle.Render(helpText)),
	)

	final := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)

	v := tea.NewView(final)
	v.AltScreen = true
	return v
}

func (m TaskModel) getQuitHelpText() string {
	switch m.InterfaceType {
	case models.InterfaceTypeWeb:
		return "Press 'q' in the server terminal to skip the task"
	case models.InterfaceTypeTUI:
		return "Press 'q' to skip task during the survey"
	case models.InterfaceTypeCLI, models.InterfaceTypeREPL:
		return `Write 'exit' to skip task during the survey

Note: Use 'tab' to scroll through suggestions`
	default:
		return m.keys.Quit.Help().Key + " to quit"
	}
}

func (m TaskModel) interfaceHelpText() string {
	switch m.InterfaceType {
	case models.InterfaceTypeWeb:
		return "Press the button in the upper right corner to view task progress"
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
