package main

import (
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/style"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const stages = 2

const (
	IntroTitle = "✦ Project Management Interface Survey ✦"

	IntroductionText = `Welcome! This survey gathers feedback on task management interfaces.

Your responses will help us improve our services. All data is anonymized
and used solely for research purposes.

By participating, you consent to data collection as described below.`

	Disclaimer = `📋 Data Collection Notice

• All responses are completely anonymized
• Data is used for research purposes only
• No personally identifiable information is collected
• You may exit at any time by pressing Esc
• We get no data unless you complete the survey and submit.`
)

type keyMap struct {
	Start    key.Binding
	Continue key.Binding
	Back     key.Binding
	Quit     key.Binding
}

var keys = keyMap{
	Start: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "start survey"),
	),
	Continue: key.NewBinding(
		key.WithKeys(" ", "j", "l", "down", "right"),
		key.WithHelp("space", "continue"),
	),
	Back: key.NewBinding(
		key.WithKeys("b", "k", "h", "backspace", "up", "left"),
		key.WithHelp("b", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c", "q"),
		key.WithHelp("esc", "quit"),
	),
}

type introModel struct {
	stage         int
	width, height int
	userQuit      bool
}

func newIntroModel() introModel {
	return introModel{
		stage: 1,
	}
}

func (m introModel) Run() error {
	model, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(introModel); ok && m.userQuit {
		return task.ErrUserQuit
	}
	return nil
}

func (m introModel) Init() tea.Cmd {
	return nil
}

func (m introModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Start) && m.stage == stages:
			return m, tea.Quit
		case key.Matches(msg, keys.Continue):
			if m.stage < stages {
				m.stage++
			}
		case key.Matches(msg, keys.Back):
			if m.stage > 1 {
				m.stage--
			}
			return m, nil
		case key.Matches(msg, keys.Quit):
			m.userQuit = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m introModel) View() string {
	if m.width < 55 || m.height < 16 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			style.TextStyle.Render("Terminal too small."))
	}

	var content string
	switch m.stage {
	case 1:
		content = IntroductionText
	case 2:
		content = Disclaimer
	default:
		return ""
	}

	boxWidth := min(m.width-10, 80)

	boxStyle := style.BorderStyle.
		Margin(1, 0).Padding(2, 4).Width(boxWidth)

	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(IntroTitle))
	b.WriteString("\n")

	b.WriteString(boxStyle.Render(style.TextStyle.Render(content)))
	b.WriteString("\n")

	helpText := "Press " + keys.Continue.Help().Key + " to continue • " +
		keys.Back.Help().Key + " to go back • " + keys.Quit.Help().Key + " to quit"

	if m.stage == stages {
		helpText += "\nPress " + keys.Start.Help().Key + " to start the survey"
	}

	b.WriteString(style.HelpStyle.Render(helpText))

	final := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String())

	return final
}

func (m *introModel) SetSize(width, height int) {
	m.width, m.height = width, height
}
