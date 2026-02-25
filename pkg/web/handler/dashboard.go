package handler

import (
	"net/http"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	app := App(r)
	hx := HTMX(r)

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	filter := models.IssueFilter{Limit: 100}
	issues, err := app.Issues.SearchIssues(r.Context(), query, filter)
	if err != nil {
		http.Error(w, "Failed to retrieve issues", http.StatusInternalServerError)
		return
	}

	var selectedIssue *models.Issue
	if len(issues) > 0 {
		selectedIssue = issues[0]
	}

	selectedIssueID := r.URL.Query().Get("selected-issue")
	if selectedIssueID != "" {
		for _, issue := range issues {
			if issue.ID == selectedIssueID {
				selectedIssue = issue
				break
			}
		}
	}

	isPartial := r.URL.Query().Get("partial") == "true"
	if hx.IsHxRequest() && isPartial {
		routes.DashboardIssueList(issues, selectedIssueID).Render(r.Context(), w)
		return
	}

	props := routes.DashboardProps{
		BaseURL:       "/",
		QueryParam:    query,
		Issues:        issues,
		SelectedIssue: selectedIssue,
		EmptyText:     "No issues found",
	}

	if hx.IsHxRequest() {
		routes.DashboardContent(props).Render(r.Context(), w)
		return
	}

	routes.Dashboard(props).Render(r.Context(), w)
}
