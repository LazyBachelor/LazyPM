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
	lastSize     tea.WindowSizeMsg
	hasSize      bool
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
	// switch dashboards based on msg
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		// remember the latest window size and forward it to the current view.
		r.hasSize = true
		r.lastSize = m
		var cmd tea.Cmd
		r.currentView, cmd = r.currentView.Update(msg)
		return r, cmd
	case msgs.SwitchToDashboardMsg:
		// switch back to dashboard 1 and apply the last known size.
		r.currentView = dashboard.NewDashboard(r.app, r.feedbackChan, r.quitChan)
		var cmds []tea.Cmd
		if r.hasSize {
			// check if there is a size, and then update it
			var sizeCmd tea.Cmd
			// set size of dashboard1 to the size before switching
			r.currentView, sizeCmd = r.currentView.Update(r.lastSize)
			if sizeCmd != nil {
				cmds = append(cmds, sizeCmd)
			}
		}
		cmds = append(cmds, tea.ClearScreen, r.currentView.Init())
		return r, tea.Batch(cmds...)
	case msgs.SwitchToDashboard2Msg:
		// switch to dashboard 2 and apply the last known size.
		r.currentView = dashboard2.NewDashboard(r.app, r.feedbackChan, r.quitChan)
		var cmds []tea.Cmd
		if r.hasSize {
			// check if there is a size, and then update it
			var sizeCmd tea.Cmd
			// set size of dashboard2 to the size before switching
			r.currentView, sizeCmd = r.currentView.Update(r.lastSize)
			if sizeCmd != nil {
				cmds = append(cmds, sizeCmd)
			}
		}
		cmds = append(cmds, tea.ClearScreen, r.currentView.Init())
		return r, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	r.currentView, cmd = r.currentView.Update(msg)
	return r, cmd
}

func (r *RootModel) View() string {
	return r.currentView.View()
}
