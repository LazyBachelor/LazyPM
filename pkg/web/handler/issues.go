package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
)

func IssuesRoutes(svc *service.Services) []Route {
	return []Route{
		{Pattern: "/issues", Handler: GetAllIssues(svc)},
		{Pattern: "POST /create-issue", Handler: CreateIssue(svc)},
	}
}

func CreateIssue(svc *service.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 1. Parse Form instead of JSON
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// 2. Map form values to your struct manually
		// (Or use a library like 'gorilla/schema')
		issue := models.Issue{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Status:      models.Status(r.FormValue("status")),
			IssueType:   models.IssueType(r.FormValue("issue_type")),
		}

		err := svc.Beads.CreateIssue(r.Context(), &issue, "")
		if err != nil {
			http.Error(w, "Failed to create issue: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. HTMX usually expects HTML back, not JSON
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<p>Created issue: %s</p>", issue.Title)
	}
}

func GetAllIssues(svc *service.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		issues, err := svc.Beads.AllIssues(r.Context())

		if err != nil {
			http.Error(w, "Failed to retrieve issues", http.StatusInternalServerError)
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
