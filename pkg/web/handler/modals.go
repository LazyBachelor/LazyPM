package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web/components"
)

func CreateIssueFormModal(w http.ResponseWriter, r *http.Request) {
	modalContent := components.IssueForm(components.IssueFormProps{
		PostAction: "/issues",
		Status:     "open",
		IssueType:  "task",
		Priority:   0,
	})

	modal := components.Modal(components.ModalProps{
		ID:      "create-issue-modal",
		Title:   "Create New Issue",
		Content: modalContent,
		Open:    true,
	})

	modal.Render(r.Context(), w)
}

func EditIssueFormModal(w http.ResponseWriter, r *http.Request) {
	issue := r.Context().Value(issueKey).(*models.Issue)

	if issue == nil {
		http.Error(w, "Issue not found in context", http.StatusInternalServerError)
		return
	}

	modalContent := components.IssueForm(components.IssueFormProps{
		PatchAction: "/issues/" + issue.ID,
		Title:       issue.Title,
		Description: issue.Description,
		Status:      string(issue.Status),
		IssueType:   string(issue.IssueType),
		Priority:    issue.Priority,
	})

	modal := components.Modal(components.ModalProps{
		ID:      "edit-issue-modal",
		Title:   "Edit Issue",
		Content: modalContent,
		Open:    true,
	})
	modal.Render(r.Context(), w)
}

func AssigneeFormModal(w http.ResponseWriter, r *http.Request) {
	issue := r.Context().Value(issueKey).(*models.Issue)

	if issue == nil {
		http.Error(w, "Issue not found in context", http.StatusInternalServerError)
		return
	}

	modalContent := components.AssigneeForm(components.AssigneeFormProps{
		PatchAction: "/issues/" + issue.ID + "/assignee",
		Assignee:    issue.Assignee,
	})

	modal := components.Modal(components.ModalProps{
		ID:      "assignee-modal",
		Title:   "Change Assignee",
		Content: modalContent,
		Open:    true,
	})
	modal.Render(r.Context(), w)
}
