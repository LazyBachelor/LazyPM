package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web/components"
	"github.com/a-h/templ"
)

type ValidationFeedback = models.ValidationFeedback

var taskFeedback ValidationFeedback
var submitChan chan<- struct{}

func SetTaskFeedback(feedback ValidationFeedback) {
	taskFeedback = feedback
}

func SetSubmitChan(ch chan<- struct{}) {
	submitChan = ch
}

func HandleTaskStatus(w http.ResponseWriter, r *http.Request) {
	if submitChan != nil {
		select {
		case submitChan <- struct{}{}:
		default:
		}
	}

	hx := HTMX(r)
	if hx.IsHxRequest() {
		hx.WriteString(`<a hx-get="/status/modal" hx-target="#modal-container" hx-swap="innerHTML">` + taskFeedback.Message + `</a>`)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	hx.WriteJSON(taskFeedback)
}

func HandleTaskStatusModal(w http.ResponseWriter, r *http.Request) {
	err := components.Modal(components.ModalProps{
		ID:      "task-status-modal",
		Title:   "Task Status",
		Content: feedbackList(taskFeedback),
		Open:    true,
	}).Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func feedbackList(feedback ValidationFeedback) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, check := range feedback.Checks {
			if !check.Valid {
				io.WriteString(w, `<p class="my-2 text-sm">`+"❌ "+check.Message+`</p>`)
			} else {
				io.WriteString(w, `<p class="my-2 text-sm">`+"✅ "+check.Message+`</p>`)
			}

		}
		return nil
	})
}
