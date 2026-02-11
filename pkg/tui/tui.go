package tui

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views"
	tea "github.com/charmbracelet/bubbletea"
)

type TUIConfig = service.Config

func Run(ctx context.Context, config TUIConfig) (tea.Model, error) {
	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return nil, err
	}

	defer cleanup()

	app := tea.NewProgram(views.NewDashboardView(svc),
		tea.WithAltScreen(), tea.WithMouseAllMotion())

	return app.Run()
}
