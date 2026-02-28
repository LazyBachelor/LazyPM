package views

import (
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views/dashboard"
)

func NewDashboardView(app *service.App, feedbackChan chan models.ValidationFeedback, quitChan chan bool) *dashboard.Model {
	return dashboard.NewDashboard(app, feedbackChan, quitChan)
}
