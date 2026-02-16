package task

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (t *Task) IntroduceTask() error {
	if t.aboutScreen == nil {
		return fmt.Errorf("aboutScreen is not set")
	}
	model, err := tea.NewProgram(t.aboutScreen, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(interface{ GetUserQuit() bool }); ok && m.GetUserQuit() {
		return ErrUserQuit
	}
	return nil
}

func (t *Task) StartInterface(ctx context.Context, cfg TaskConfig) error {
	if t.interfaceType == nil {
		return fmt.Errorf("interfaceType is not set")
	}

	return t.interfaceType.Run(ctx, cfg)
}

func (t *Task) StartQuestionnaire() error {
	if t.questionnaire == nil {
		return fmt.Errorf("questionnaire is not set")
	}
	model, err := tea.NewProgram(t.questionnaire, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}
	if m, ok := model.(interface{ GetUserQuit() bool }); ok && m.GetUserQuit() {
		return ErrUserQuit
	}
	return nil
}
