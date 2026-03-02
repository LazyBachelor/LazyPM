package views

import (
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views/dashboard"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views/dashboard2"
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	currentView  tea.Model
	app          *app.App
	feedbackChan chan models.ValidationFeedback
	quitChan     chan bool
}

func NewRootView(app *app.App, feedbackChan chan models.ValidationFeedback, quitChan chan bool) *RootModel {
	initialView := dashboard.NewDashboard(app, feedbackChan, quitChan)
	return &RootModel{
		currentView:  initialView,
		app:          app,
		feedbackChan: feedbackChan,
		quitChan:     quitChan,
	}
}

func (r *RootModel) Init() tea.Cmd {
	return r.currentView.Init()
}

func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case msgs.SwitchToDashboard2Msg:
		r.currentView = dashboard2.NewDashboard(r.app, r.feedbackChan, r.quitChan)
		return r, r.currentView.Init()
	case msgs.SwitchToDashboardMsg:
		r.currentView = dashboard.NewDashboard(r.app, r.feedbackChan, r.quitChan)
		return r, r.currentView.Init()
	}

	var cmd tea.Cmd
	r.currentView, cmd = r.currentView.Update(msg)
	return r, cmd
}

func (r *RootModel) View() string {
	return r.currentView.View()
}
