package handler

import (
	"encoding/json"
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

		var issue models.Issue

		err := json.NewDecoder(r.Body).Decode(&issue)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		err = svc.Beads.CreateIssue(r.Context(), &issue, "")
		if err != nil {
			http.Error(w, "Failed to create issue: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(issue)
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
