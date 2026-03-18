package tui

import (
	"context"

	"charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views"
)

type Config = models.Config
type ValidationFeedback = models.ValidationFeedback

type Tui struct {
	feedbackChan chan ValidationFeedback
	quitChan     chan bool
	submitChan   chan<- struct{}
}

func New() *Tui {
	return &Tui{}
}

func (t *Tui) Run(ctx context.Context, config Config) error {
	app, cleanup, err := app.New(ctx, config)
	if err != nil {
		return err
	}

	defer cleanup()

	p := tea.NewProgram(views.NewRootView(app, t.feedbackChan, t.quitChan, t.submitChan),
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

func (t *Tui) SetChannels(feedbackChan chan ValidationFeedback, quitChan chan bool) {
	t.feedbackChan = feedbackChan
	t.quitChan = quitChan
}

func (t *Tui) SetSubmitChan(submitChan chan<- struct{}) {
	t.submitChan = submitChan
}
