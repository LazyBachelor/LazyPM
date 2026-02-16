package main

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

var ErrUserQuit = errors.New("user quit")

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
	stage         int
	width, height int
	userQuit      bool
}

func newIntroModel() introModel {
	return introModel{
		stage: 0,
	}
}

func (m introModel) Init() tea.Cmd {
	return nil
}

func (m introModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeySpace:
			m.stage++
			if m.stage > 1 {
				return m, tea.Quit
			}
		case tea.KeyEsc, tea.KeyCtrlC:
			m.userQuit = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m introModel) View() string {
	switch m.stage {
	case 0:
		return IntroductionText
	case 1:
		return Disclaimer
	}
	return ""
}

func (m introModel) Run() error {
	model, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(introModel); ok && m.userQuit {
		return ErrUserQuit
	}
	return nil
}

func (m *introModel) SetSize(width, height int) {
	m.width, m.height = width, height
}
