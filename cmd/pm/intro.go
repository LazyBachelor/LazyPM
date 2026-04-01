package main

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/style"
)

const stages = 2

const (
	IntroTitle       = "Project Management Interface Survey"
	IntroductionText = `This survey evaluates how people use different project management interfaces to complete common tasks.
You will complete a short set of task-based exercises and answer a few follow-up questions about the experience. Your feedback helps us compare interface design and usability.

The survey takes around 15-25 minutes to complete.`

	Disclaimer = `Data Collection Notice

• Responses are anonymized before analysis
• We do not collect personally identifying information
• You may exit at any time by pressing Ctrl+C between tasks or closing the terminal.
• We collect task results, questionnaire answers, and interaction events during the survey
• Data is collected continuously during the survey, so exiting early will still contribute valuable insights.`
)

type keyMap struct {
	Start    key.Binding
	Continue key.Binding
	Back     key.Binding
	Quit     key.Binding
}

type introModel struct {
	stage         int
	width, height int
	userQuit      bool
	keys          keyMap
}

func newIntroModel() introModel {
	var keys = keyMap{
		Start: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "start survey"),
		),
		Continue: key.NewBinding(
			key.WithKeys("space", "j", "l", "down", "right"),
			key.WithHelp("space", "continue"),
		),
		Back: key.NewBinding(
			key.WithKeys("b", "k", "h", "backspace", "up", "left"),
			key.WithHelp("b", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "quit"),
		),
	}
	return introModel{
		stage: 1,
		keys:  keys,
	}
}

func (m introModel) Run() error {
	model, err := tea.NewProgram(m).Run()
	if err != nil {
		return err
	}
	if im, ok := model.(introModel); ok && im.userQuit {
		return models.ErrUserQuit
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
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Start) && m.stage == stages:
			return m, tea.Quit
		case key.Matches(msg, m.keys.Continue):
			if m.stage < stages {
				m.stage++
			}
		case key.Matches(msg, m.keys.Back):
			if m.stage > 1 {
				m.stage--
			}
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.userQuit = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m introModel) View() tea.View {
	if m.width < 55 || m.height < 16 {
		content := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			style.TextStyle.Render("Terminal too small."))
		v := tea.NewView(content)
		v.AltScreen = true
		return v
	}

	var content string
	switch m.stage {
	case 1:
		content = IntroductionText
	case 2:
		content = Disclaimer
	default:
		return tea.NewView("")
	}

	boxWidth := min(m.width-10, 120)

	boxStyle := style.BorderStyle.
		Margin(1, 0).Padding(2, 4).Width(boxWidth)

	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(IntroTitle))
	b.WriteString("\n")

	b.WriteString(boxStyle.Render(style.TextStyle.Render(content)))
	b.WriteString("\n")

	helpText := "Press " + m.keys.Continue.Help().Key + " to continue • " +
		m.keys.Back.Help().Key + " to go back • " + m.keys.Quit.Help().Key + " to quit"

	if m.stage == stages {
		helpText += "\nPress " + m.keys.Start.Help().Key + " to start the survey"
	}

	b.WriteString(style.HelpStyle.Render(helpText))

	final := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String())

	v := tea.NewView(final)
	v.AltScreen = true
	return v
}

func (m *introModel) SetSize(width, height int) {
	m.width, m.height = width, height
}
