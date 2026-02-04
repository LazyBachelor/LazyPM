package views

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views/dashboard"
)

func NewDashboardView(svc *service.Services) dashboard.Model {
	return dashboard.NewDashboard(svc)
}
