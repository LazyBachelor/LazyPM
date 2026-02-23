package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web/routes"
	"github.com/go-chi/chi/v5"
)

const issueKey = "issue"
const commentsKey = "comments"

type IssueForm struct {
	Title       string           `form:"title" validate:"required,max=255"`
	Description string           `form:"description" validate:"required,max=2000"`
	Status      models.Status    `form:"status" validate:"required,oneof=open in_progress closed"`
	IssueType   models.IssueType `form:"issue_type" validate:"required,oneof=task bug feature chore"`
	Priority    int              `form:"priority" validate:"gte=0,lte=4"`
}

type UpdateIssueForm struct {
	Title       *string           `form:"title" validate:"omitempty,max=255"`
	Description *string           `form:"description" validate:"omitempty,max=2000"`
	Status      *models.Status    `form:"status" validate:"omitempty,oneof=open in_progress closed"`
	IssueType   *models.IssueType `form:"issue_type" validate:"omitempty,oneof=task bug feature chore"`
	Priority    *int              `form:"priority" validate:"omitempty,gte=0,lte=4"`
}

func CreateIssue(w http.ResponseWriter, r *http.Request) {
	app := App(r)
	hx := HTMX(r)

	form, err := ParseForm[IssueForm](r)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	if err := ValidateForm(form); err != nil {
		if hx.IsHxRequest() {
			hx.WriteString("<div>Please fix the form errors</div>")
		} else {
			http.Error(w, "Validation error: "+err.Error(), http.StatusUnprocessableEntity)
		}
		return
	}

	issue := form.toIssue()
	if err := app.Issues.CreateIssue(r.Context(), &issue, ""); err != nil {
		http.Error(w, "Failed to create issue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if hx.IsHxRequest() {
		hx.WriteString("<div>Issue created successfully</div>")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	hx.WriteJSON(issue)
}

func ListIssues(w http.ResponseWriter, r *http.Request) {
	app := App(r)
	hx := HTMX(r)

	issues, err := app.Issues.AllIssues(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve issues", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	hx.WriteJSON(issues)
}

func IssueCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app := App(r)

		id := chi.URLParam(r, "id")
		issue, err := app.Issues.GetIssue(r.Context(), id)
		if err != nil {
			http.Error(w, "Error getting issue: "+err.Error(), http.StatusNotFound)
			return
		}
		if issue == nil {
			http.Error(w, "Issue not found", http.StatusNotFound)
			return
		}

		comments, err := app.Issues.GetIssueComments(r.Context(), issue.ID)
		if err != nil {
			http.Error(w, "Error getting comments: "+err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), issueKey, issue)
		ctx = context.WithValue(ctx, commentsKey, comments)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func GetIssue(w http.ResponseWriter, r *http.Request) {
	issue := r.Context().Value(issueKey).(*models.Issue)
	comments := r.Context().Value(commentsKey).([]*models.Comment)
	hx := HTMX(r)

	if !hx.IsHxRequest() && strings.Contains(r.Header.Get("Accept"), "text/html") {
		routes.IssueDetail(routes.IssueDetailProps{
			Issue:    *issue,
			Comments: comments,
		}).Render(r.Context(), w)
		return
	}

	if hx.IsHxRequest() {
		routes.IssueDetailContent(routes.IssueDetailProps{
			Issue:    *issue,
			Comments: comments,
		}).Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	hx.WriteJSON(issue)
}

func UpdateIssue(w http.ResponseWriter, r *http.Request) {
	issue := r.Context().Value(issueKey).(*models.Issue)
	app := App(r)
	hx := HTMX(r)

	form, err := ParseForm[UpdateIssueForm](r)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	if err := ValidateForm(form); err != nil {
		if hx.IsHxRequest() {
			hx.WriteString("<div>Please fix the form errors</div>")
		} else {
			http.Error(w, "Validation error: "+err.Error(), http.StatusUnprocessableEntity)
		}
		return
	}

	changes := form.toChanges()

	if err := app.Issues.UpdateIssue(r.Context(), issue.ID, changes, ""); err != nil {
		http.Error(w, "Failed to update issue", http.StatusInternalServerError)
		return
	}

	issue, err = app.Issues.GetIssue(r.Context(), issue.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve updated issue", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	hx.WriteJSON(issue)
}

func DeleteIssue(w http.ResponseWriter, r *http.Request) {
	issue := r.Context().Value(issueKey).(*models.Issue)

	app := App(r)
	if err := app.Issues.DeleteIssue(r.Context(), issue.ID); err != nil {
		http.Error(w, "Failed to delete issue", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (f *IssueForm) toIssue() models.Issue {
	return models.Issue{
		Title:       f.Title,
		Description: f.Description,
		Status:      f.Status,
		IssueType:   f.IssueType,
		Priority:    f.Priority,
	}
}

func (f *UpdateIssueForm) toChanges() map[string]any {
	changes := make(map[string]any)
	if f.Title != nil {
		changes["title"] = *f.Title
	}
	if f.Description != nil {
		changes["description"] = *f.Description
	}
	if f.Status != nil {
		changes["status"] = *f.Status
	}
	if f.IssueType != nil {
		changes["issue_type"] = *f.IssueType
	}
	if f.Priority != nil {
		changes["priority"] = *f.Priority
	}
	return changes
}
