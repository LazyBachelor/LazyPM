package handler

import (
	"context"
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/go-chi/chi/v5"
)

type IssueForm struct {
	Title       string           `form:"title" validate:"required,max=255"`
	Description string           `form:"description" validate:"required,max=2000"`
	Status      models.Status    `form:"status" validate:"required,oneof=open in_progress closed"`
	IssueType   models.IssueType `form:"issue_type" validate:"required,oneof=task bug feature chore"`
	Priority    int              `form:"priority" validate:"gte=0,lte=4"`
}

func (f *IssueForm) ToIssue() models.Issue {
	return models.Issue{
		Title:       f.Title,
		Description: f.Description,
		Status:      f.Status,
		IssueType:   f.IssueType,
		Priority:    f.Priority,
	}
}

func CreateIssue(w http.ResponseWriter, r *http.Request) {
	svc := Services(r)
	hx := HTMX(r)

	form, err := ParseForm[IssueForm](r)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	if err := ValidateForm(form); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if hx.IsHxRequest() {
			hx.WriteString("<div class='alert alert-error'>Please fix the form errors</div>")
		} else {
			hx.WriteJSON(map[string]interface{}{"error": err.Error()})
		}
		return
	}

	issue := form.ToIssue()
	if err := svc.Beads.CreateIssue(r.Context(), &issue, ""); err != nil {
		http.Error(w, "Failed to create issue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if hx.IsHxRequest() {
		hx.WriteString("<div>Issue created successfully</div>")
		return
	}

	hx.WriteJSON(map[string]any{
		"title":  issue.Title,
		"status": issue.Status,
	})
}

func ListIssues(w http.ResponseWriter, r *http.Request) {
	svc := Services(r)
	hx := HTMX(r)

	issues, err := svc.Beads.AllIssues(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve issues", http.StatusInternalServerError)
		return
	}

	hx.WriteJSON(issues)
}

func IssueCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		svc := Services(r)

		id := chi.URLParam(r, "id")
		issue, err := svc.Beads.GetIssue(r.Context(), id)
		if err != nil {
			http.Error(w, "Issue not found", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "issue", issue)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func GetIssue(w http.ResponseWriter, r *http.Request) {
	issue := r.Context().Value("issue").(*models.Issue)
	hx := HTMX(r)

	hx.WriteJSON(issue)
}

func UpdateIssue(w http.ResponseWriter, r *http.Request) {
	issue := r.Context().Value("issue").(*models.Issue)
	svc := Services(r)
	hx := HTMX(r)

	changes := make(map[string]any)

	if err := svc.Beads.UpdateIssue(r.Context(), issue.ID, changes, ""); err != nil {
		http.Error(w, "Failed to update issue", http.StatusInternalServerError)
		return
	}

	issue, err := svc.Beads.GetIssue(r.Context(), issue.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve updated issue", http.StatusInternalServerError)
		return
	}

	hx.WriteJSON(issue)
}

func DeleteIssue(w http.ResponseWriter, r *http.Request) {
	issue := r.Context().Value("issue").(*models.Issue)

	svc := Services(r)
	if err := svc.Beads.DeleteIssue(r.Context(), issue.ID); err != nil {
		http.Error(w, "Failed to delete issue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
