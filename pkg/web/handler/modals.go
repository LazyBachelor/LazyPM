package handler

import (
	"net/http"
	"net/url"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web/components"
)

func CreateIssueFormModal(w http.ResponseWriter, r *http.Request) {
	// Get the 'from' parameter to track where the user came from
	from := r.URL.Query().Get("from")
	postAction := "/issues"
	if from != "" {
		v := url.Values{}
		v.Set("from", from)
		postAction += "?" + v.Encode()
	}

	modalContent := components.IssueForm(components.IssueFormProps{
		PostAction: postAction,
		Status:     "open",
		IssueType:  "task",
		Priority:   0,
		Assignee:   "",
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
	issueVal := r.Context().Value(issueKey)
	issue, ok := issueVal.(*models.Issue)
	if !ok || issue == nil {
		http.Error(w, "Issue not found in context", http.StatusInternalServerError)
		return
	}

	// Get the 'from' parameter to track where the user came from
	from := r.URL.Query().Get("from")
	patchAction := "/issues/" + issue.ID
	deleteAction := "/issues/" + issue.ID + "/delete"
	if from != "" {
		v := url.Values{}
		v.Set("from", from)
		q := v.Encode()
		patchAction += "?" + q
		deleteAction += "?" + q
	}

	modalContent := components.IssueForm(components.IssueFormProps{
		PatchAction:  patchAction,
		DeleteAction: deleteAction,
		IssueID:      issue.ID,
		Title:        issue.Title,
		Description:  issue.Description,
		Status:       string(issue.Status),
		CloseReason:  issue.CloseReason,
		IssueType:    string(issue.IssueType),
		Priority:     issue.Priority,
		Assignee:     issue.Assignee,
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
	issueVal := r.Context().Value(issueKey)
	issue, ok := issueVal.(*models.Issue)
	if !ok || issue == nil {
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
		v := url.Values{}
		v.Set("from", from)
		deleteAction += "?" + v.Encode()
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

func CloseIssueFormModal(w http.ResponseWriter, r *http.Request) {
	issueVal := r.Context().Value(issueKey)
	issue, ok := issueVal.(*models.Issue)
	if !ok || issue == nil {
		http.Error(w, "Issue not found in context", http.StatusInternalServerError)
		return
	}

	if issue.Status == models.StatusClosed {
		http.Error(w, "Issue is already closed", http.StatusBadRequest)
		return
	}

	modalContent := components.CloseIssueForm(components.CloseIssueFormProps{
		PostAction: "/issues/" + issue.ID + "/close",
	})

	modal := components.Modal(components.ModalProps{
		ID:      "close-issue-modal",
		Title:   "Close Issue",
		Content: modalContent,
		Open:    true,
	})
	modal.Render(r.Context(), w)
}
