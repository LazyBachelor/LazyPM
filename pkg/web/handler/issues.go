package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
)

func IssuesRoutes(svc *service.Services) []Route {
	return []Route{
		{Pattern: "/api/issues", Handler: GetAllIssues(svc)},
		{Pattern: "POST /create-issue", Handler: CreateIssue(svc)},
	}
}

func CreateIssue(svc *service.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		priority, _ := strconv.Atoi(r.FormValue("priority"))
		issue := models.Issue{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Status:      models.Status(r.FormValue("status")),
			IssueType:   models.IssueType(r.FormValue("issue_type")),
			Priority:    priority,
		}

		err := svc.Beads.CreateIssue(r.Context(), &issue, "")
		if err != nil {
			http.Error(w, "Failed to create issue: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/issues", http.StatusSeeOther)
	}
}

func GetAllIssues(svc *service.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		issues, err := svc.Beads.AllIssues(r.Context())
		if err != nil {
			http.Error(w, errRetrieveIssues, http.StatusInternalServerError)
			return
		}
		jsonData, err := json.Marshal(issues)
		if err != nil {
			http.Error(w, "Failed to marshal issues", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func IssuesHandler(svc *service.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/issues" {
			handleNotFound(w, r)
			return
		}

		issues, err := svc.Beads.AllIssues(r.Context())
		if err != nil {
			http.Error(w, errRetrieveIssues, http.StatusInternalServerError)
			return
		}

		openCount, closedCount := models.CountByStatus(issues)
		filterPri := r.URL.Query().Get("pri")
		filterState := r.URL.Query().Get("state")
		filtered := models.FilterIssues(issues, filterPri, filterState)
		models.SortIssuesByPriority(filtered, true)

		props := routes.IssuesProps{
			Issues:      filtered,
			FilterPri:   filterPri,
			FilterState: filterState,
			OpenCount:   openCount,
			ClosedCount: closedCount,
		}
		routes.Issues(props).Render(r.Context(), w)
	}
}

func IssueDetailHandler(svc *service.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/issues/")
		if id == "" {
			http.Redirect(w, r, "/issues", http.StatusFound)
			return
		}

		issues, err := svc.Beads.AllIssues(r.Context())
		if err != nil {
			http.Error(w, errRetrieveIssues, http.StatusInternalServerError)
			return
		}

		found := models.FindIssueByID(issues, id)
		if found == nil {
			handleNotFound(w, r)
			return
		}

		routes.IssueDetail(routes.IssueDetailProps{Issue: *found}).Render(r.Context(), w)
	}
}
