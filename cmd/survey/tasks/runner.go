package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

func (t *Task) IntroduceTask() error {
	_, err := tea.NewProgram(t.aboutScreen, tea.WithAltScreen()).Run()
	return err
}

func (t *Task) StartInterface(ctx context.Context, cfg service.Config) error {
	return t.interfaceType.Run(ctx, cfg)
}

func (t *Task) StartQuestionnaire() error {
	_, err := tea.NewProgram(t.questionnaire, tea.WithAltScreen()).Run()
	return err
}
