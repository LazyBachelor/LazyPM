package handler

import (
	"fmt"
	"net/http"
	"sort"
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

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Priority > issues[j].Priority
	})

	isBoardView := r.URL.Query().Get("board") == "true" || strings.HasPrefix(r.URL.Path, "/board/")

	if isBoardView {
		ctx := r.Context()

		backlogIssues, err := app.Issues.GetIssuesNotInAnySprint(ctx)
		if err != nil {
			backlogIssues = []*models.Issue{}
		}

		sort.Slice(backlogIssues, func(i, j int) bool {
			return backlogIssues[i].Priority > backlogIssues[j].Priority
		})

		sprints, err := app.Issues.GetSprints(ctx)
		if err != nil {
			sprints = []int{}
		}

		currentSprint := 0
		sprintParam := r.URL.Query().Get("sprint")
		if sprintParam != "" {
			fmt.Sscanf(sprintParam, "%d", &currentSprint)
		}

		if currentSprint == 0 && len(sprints) > 0 {
			currentSprint = sprints[0]
		}

		var sprintIssues []*models.Issue
		if currentSprint > 0 {
			sprintIssues, err = app.Issues.GetIssuesBySprint(ctx, currentSprint)
			if err != nil {
				sprintIssues = []*models.Issue{}
			}
		}

		sort.Slice(sprintIssues, func(i, j int) bool {
			return sprintIssues[i].Priority > sprintIssues[j].Priority
		})

		boardProps := routes.BoardViewProps{
			BaseURL:       "/?board=true",
			QueryParam:    query,
			Issues:        issues,
			BacklogIssues: backlogIssues,
			SprintIssues:  sprintIssues,
			CurrentSprint: currentSprint,
			Sprints:       sprints,
			EmptyText:     "No issues found",
		}

		if hx.IsHxRequest() {
			if r.URL.Query().Get("board") != "true" {
				w.Header().Set("HX-Push-Url", fmt.Sprintf("/?board=true&sprint=%d", currentSprint))
			}

			if strings.HasPrefix(r.URL.Path, "/board/") {
				routes.BoardColumns(boardProps).Render(r.Context(), w)
			} else {
				routes.BoardViewContent(boardProps).Render(r.Context(), w)
			}
			return
		}

		routes.BoardView(boardProps).Render(r.Context(), w)
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

	if selectedIssue != nil {
		comments, err := app.Issues.GetIssueComments(r.Context(), selectedIssue.ID)
		if err == nil && comments != nil {
			selectedIssue.Comments = comments
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

func CreateSprintHandler(w http.ResponseWriter, r *http.Request) {
	app := App(r)
	ctx := r.Context()

	sprintNum, err := app.Issues.AddSprint(ctx)
	if err != nil {
		http.Error(w, "Failed to create sprint", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("/?board=true&sprint=%d", sprintNum)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
}
