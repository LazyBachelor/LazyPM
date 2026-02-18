package tui

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views"
	tea "github.com/charmbracelet/bubbletea"
)

type TUIConfig = service.Config

type Tui struct {
	feedbackChan chan task.ValidationFeedback
	quitChan     chan bool
}

func NewTui() *Tui {
	return &Tui{}
}

func (t *Tui) Run(ctx context.Context, config TUIConfig) error {
	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	p := tea.NewProgram(views.NewDashboardView(svc, t.feedbackChan, t.quitChan),
		tea.WithAltScreen(), tea.WithMouseAllMotion())

	if t.quitChan != nil {
		go func() {
			<-t.quitChan
			p.Quit()
		}()
	}

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func (t *Tui) SetChannels(feedbackChan chan task.ValidationFeedback, quitChan chan bool) {
	t.feedbackChan = feedbackChan
	t.quitChan = quitChan
}
