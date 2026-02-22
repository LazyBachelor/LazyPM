package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web/components"
)

type CommentForm struct {
	Author string `form:"author" validate:"required,max=100"`
	Text   string `form:"text" validate:"required,max=2000"`
}

func ListComments(w http.ResponseWriter, r *http.Request) {
	comments := r.Context().Value(commentsKey).([]*models.Comment)
	hx := HTMX(r)

	if hx.IsHxRequest() {
		components.CommentList(components.CommentListProps{
			Comments: comments,
		}).Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	hx.WriteJSON(comments)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	issue := r.Context().Value(issueKey).(*models.Issue)
	app := App(r)
	hx := HTMX(r)

	form, err := ParseForm[CommentForm](r)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	if err := ValidateForm(form); err != nil {
		if hx.IsHxRequest() {
			hx.WriteString(`<div class="alert alert-error">Please fix the form errors</div>`)
		} else {
			http.Error(w, "Validation error: "+err.Error(), http.StatusUnprocessableEntity)
		}
		return
	}

	comment, err := app.Issues.AddIssueComment(r.Context(), issue.ID, form.Author, form.Text)
	if err != nil {
		http.Error(w, "Failed to create comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if hx.IsHxRequest() {
		components.CommentItem(components.CommentItemProps{
			Comment: comment,
		}).Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	hx.WriteJSON(comment)
}
