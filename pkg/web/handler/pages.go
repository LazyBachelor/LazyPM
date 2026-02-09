package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/service"
)

const errRetrieveIssues = "Failed to retrieve issues"

func PagesRoutes(svc *service.Services) []Route {
	return []Route{
		{Pattern: "/", Handler: IndexHandler(svc)},
		{Pattern: "/boards", Handler: BoardsHandler(svc)},
		{Pattern: "/create", Handler: CreateHandler()},
		{Pattern: "/dashboards/new", Handler: NewDashboardHandler()},
		{Pattern: "/issues", Handler: IssuesHandler(svc)},
		{Pattern: "/issues/", Handler: IssueDetailHandler(svc)},
	}
}

func handleNotFound(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Page not found", http.StatusNotFound)
}
