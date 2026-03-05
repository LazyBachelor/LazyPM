package views

import (
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views/dashboard"
)

func NewDashboardView(app *app.App, feedbackChan chan models.ValidationFeedback, quitChan chan bool, submitChan chan<- struct{}) *dashboard.Model {
	return dashboard.NewDashboard(app, feedbackChan, quitChan, submitChan)
}
