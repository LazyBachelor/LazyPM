package handler

import (
	"net/http"
	"sort"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web/components"
	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
)

func orDisplay(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

func selectedID(issue *models.Issue) string {
	if issue == nil {
		return ""
	}
	return issue.ID
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	app := App(r)

	issues, err := app.Issues.AllIssues(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve issues", http.StatusInternalServerError)
		return
	}

	var activeIssues, closedIssues []models.Issue
	for _, issue := range issues {
		if issue.Status == models.StatusClosed {
			closedIssues = append(closedIssues, issue)
		} else {
			activeIssues = append(activeIssues, issue)
		}
	}
	sort.Slice(activeIssues, func(i, j int) bool { return activeIssues[i].Priority > activeIssues[j].Priority })
	sort.Slice(closedIssues, func(i, j int) bool { return closedIssues[i].Priority > closedIssues[j].Priority })

	selID := r.URL.Query().Get("id")
	var selectedIssue *models.Issue
	if selID != "" {
		for i := range issues {
			if issues[i].ID == selID {
				selectedIssue = &issues[i]
				break
			}
		}
	}
	if selectedIssue == nil && len(activeIssues) > 0 {
		selectedIssue = &activeIssues[0]
		selID = selectedIssue.ID
	}

	searchQuery := r.URL.Query().Get("q")

	var displayOwner, displayAssignee string
	if selectedIssue != nil {
		displayOwner = orDisplay(selectedIssue.Owner, "—")
		displayAssignee = orDisplay(selectedIssue.Assignee, "—")
	}

	props := routes.DashboardProps{
		ActiveIssues:    activeIssues,
		ClosedIssues:    closedIssues,
		SelectedIssue:   selectedIssue,
		SelectedID:      selectedID(selectedIssue),
		SearchQuery:     searchQuery,
		DisplayOwner:    displayOwner,
		DisplayAssignee: displayAssignee,
	}

	if HTMX(r).IsHxRequest() {
		components.DashboardPage(components.DashboardPageProps{
			ActiveIssues:    props.ActiveIssues,
			ClosedIssues:    props.ClosedIssues,
			SelectedIssue:   props.SelectedIssue,
			SelectedID:      props.SelectedID,
			BaseURL:         "/dashboard",
			QueryParam:      "?id=",
			EmptyText:       "Select an issue to view details",
			DisplayOwner:    props.DisplayOwner,
			DisplayAssignee: props.DisplayAssignee,
		}).Render(r.Context(), w)
		return
	}

	routes.Dashboard(props).Render(r.Context(), w)
} 