package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/pkg/web/components"
)

func CreateIssueFormModal(w http.ResponseWriter, r *http.Request) {
	modalContent := components.IssueForm(components.IssueFormProps{
		Action:    "/issues",
		Status:    "open",
		IssueType: "task",
		Priority:  0,
	})

	modal := components.Modal(components.ModalProps{
		ID:      "create-issue-modal",
		Title:   "Create New Issue",
		Content: modalContent,
		Open:    true,
	})

	modal.Render(r.Context(), w)
}
