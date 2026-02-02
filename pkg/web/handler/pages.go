package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/web/components"
	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
)

func PagesRoutes(svc *service.Services) []Route {
	return []Route{
		{Pattern: "/", Handler: IndexHandler(svc)},
	}
}

func IndexHandler(svc *service.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			handleNotFound(w, r)
			return
		}

		issues, err := svc.Beads.AllIssues(r.Context())

		if err != nil {
			http.Error(w, "Failed to retrieve issues",
				http.StatusInternalServerError)
			return
		}

		props := routes.IndexProps{
			Issues: issues,
			IssueList: components.IssueListProps{
				Issues: issues,
			},
		}

		routes.Index(props).Render(r.Context(), w)
	}
}

func handleNotFound(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Page not found", http.StatusNotFound)
}
