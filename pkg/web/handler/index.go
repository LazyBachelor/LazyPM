package handler

import (
	"net/http"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web/components"
	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	app := App(r)
	hx := HTMX(r)

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	filter := models.IssueFilter{Limit: 100}
	issuesPtr, err := app.Issues.SearchIssues(r.Context(), query, filter)
	if err != nil {
		http.Error(w, "failed to retrieve issues", http.StatusInternalServerError)
		return
	}
	issues := models.IssuesPtrToIssues(issuesPtr)

	props := routes.IndexProps{
		IssueTable: components.IssueTableProps{
			Issues: issues,
		},
		SearchQuery: query,
	}

	if hx.IsHxRequest() {
		routes.IndexContent(props).Render(r.Context(), w)
		return
	}

	routes.Index(props).Render(r.Context(), w)
}
