package tui

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views"
	tea "github.com/charmbracelet/bubbletea"
)

type TUIConfig = service.Config

type Tui struct{}

func NewTui() *Tui {
	return &Tui{}
}

func (t *Tui) Run(ctx context.Context, config TUIConfig) error {
	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return nil
	}

	defer cleanup()

	if _, err := tea.NewProgram(views.NewDashboardView(svc),
		tea.WithAltScreen(), tea.WithMouseAllMotion()).Run(); err != nil {
		return err
	}

	return nil
}
