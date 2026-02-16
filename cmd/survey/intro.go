package main

import (
	"github.com/LazyBachelor/LazyPM/internal/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

const (
	IntroTitle = "Welcome to the Task Survey!"

	IntroductionText = `This survey is designed to gather feedback on task management interfaces.
Your responses will help us improve our services. Please note that all data collected will be anonymized and used solely for research purposes.
By participating, you consent to the collection and use of your data as described in this disclaimer.

Press any key to continue...`

	Disclaimer = `This survey is designed to gather feedback on task management interfaces.
Your responses will help us improve our services. Please note that all data collected will be anonymized and used solely for research purposes.
By participating, you consent to the collection and use of your data as described in this disclaimer.`
)

type introModel struct {
	form          *huh.Form
	width, height int
}

func newIntroModel() introModel {
	return introModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewNote().Title(IntroTitle).Description(IntroductionText),
			),
			huh.NewGroup(
				huh.NewNote().Title("Disclaimer").Description(Disclaimer),
			),
		),
	}
}

func (m introModel) Init() tea.Cmd {
	m.form.Init()
	return nil
}

func (m introModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
			case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	return m, cmd
}

func (m introModel) View() string {
	return m.form.WithTheme(style.HuhCenterTheme()).View()
}

func (m introModel) Run() error {
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

func (m *introModel) SetSize(width, height int) {
	m.width, m.height = width, height
}
