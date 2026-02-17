package views

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views/dashboard"
)

func NewDashboardView(svc *service.Services, feedbackChan chan task.ValidationFeedback, quitChan chan bool) *dashboard.Model {
	return dashboard.NewDashboard(svc, feedbackChan, quitChan)
}
