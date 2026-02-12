package tasks

import (
	"context"
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

func (t *Task) IntroduceTask() error {
	if t.aboutScreen == nil {
		return fmt.Errorf("aboutScreen is not set")
	}
	_, err := tea.NewProgram(t.aboutScreen, tea.WithAltScreen()).Run()
	return err
}

func (t *Task) StartInterface(ctx context.Context, cfg service.Config) error {
	if t.interfaceType == nil {
		return fmt.Errorf("interfaceType is not set")
	}

	return t.interfaceType.Run(ctx, cfg)
}

func (t *Task) StartQuestionnaire() error {
	if t.questionnaire == nil {
		return fmt.Errorf("questionnaire is not set")
	}
	_, err := tea.NewProgram(t.questionnaire, tea.WithAltScreen()).Run()
	return err
}
