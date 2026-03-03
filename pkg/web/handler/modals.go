package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web/components"
)

func CreateIssueFormModal(w http.ResponseWriter, r *http.Request) {
	// Get the 'from' parameter to track where the user came from
	from := r.URL.Query().Get("from")
	postAction := "/issues"
	if from != "" {
		postAction += "?from=" + from
	}

	modalContent := components.IssueForm(components.IssueFormProps{
		PostAction: postAction,
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

	// Get the 'from' parameter to track where the user came from
	from := r.URL.Query().Get("from")
	patchAction := "/issues/" + issue.ID
	deleteAction := "/issues/" + issue.ID + "/delete"
	if from != "" {
		patchAction += "?from=" + from
		deleteAction += "?from=" + from
	}

	modalContent := components.IssueForm(components.IssueFormProps{
		PatchAction:  patchAction,
		DeleteAction: deleteAction,
		Title:        issue.Title,
		Description:  issue.Description,
		Status:       string(issue.Status),
		IssueType:    string(issue.IssueType),
		Priority:     issue.Priority,
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

func DeleteIssueConfirmModal(w http.ResponseWriter, r *http.Request) {
	issue := r.Context().Value(issueKey).(*models.Issue)

	if issue == nil {
		http.Error(w, "Issue not found in context", http.StatusInternalServerError)
		return
	}

	// Get the 'from' parameter to track where the user came from
	from := r.URL.Query().Get("from")
	deleteAction := "/issues/" + issue.ID
	if from != "" {
		deleteAction += "?from=" + from
	}

	confirmModal := components.ConfirmModal(components.ConfirmModalProps{
		Title:         "Delete Issue",
		Message:       "Are you sure you want to delete this issue? This action cannot be undone.",
		ConfirmText:   "Delete Issue",
		CancelText:    "Cancel",
		ConfirmAction: deleteAction,
		ConfirmClass:  "btn-error",
		IssueID:       issue.ID,
	})

	confirmModal.Render(r.Context(), w)
}
