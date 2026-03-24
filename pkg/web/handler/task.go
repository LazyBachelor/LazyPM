package handler

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/web/components"
	"github.com/a-h/templ"
)

type ValidationFeedback = models.ValidationFeedback

var taskFeedback ValidationFeedback
var submitChan chan<- models.ValidationTrigger

func SetTaskFeedback(feedback ValidationFeedback) {
	taskFeedback = feedback
}

func SetSubmitChan(ch chan<- models.ValidationTrigger) {
	submitChan = ch
}

func HandleTaskStatus(w http.ResponseWriter, r *http.Request) {
	if submitChan != nil {
		select {
		case submitChan <- models.ValidationTrigger{Source: models.ValidationTriggerAutoPoll}:
		default:
		}
	}
	time.Sleep(100 * time.Millisecond)

	hx := HTMX(r)
	if hx.IsHxRequest() {
		if taskFeedback.Success {
			w.Header().Set("HX-Trigger", "task-status-success")
			w.WriteHeader(http.StatusNoContent)
		} else {
			hx.WriteString(`
		<div id="status" type="button" hx-get="/status" hx-target="#status" hx-swap="outerHTML">
				<button class="btn btn-error btn-sm" hx-get="/status/modal" hx-target="#modal-container" hx-swap="innerHTML">` + taskFeedback.Message + `</button>
		</div>`)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	hx.WriteJSON(taskFeedback)
}

func HandleTaskStatusModal(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond)
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
		if feedback.Success {
			io.WriteString(w, `<p class="my-2">`+feedback.Message+`</p>`)
			io.WriteString(w, `<p class="my-2 text-xl">You can close the browser and return to the terminal window</p>`)
		}

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
