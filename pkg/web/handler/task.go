package handler

import (
	"net/http"

	"github.com/LazyBachelor/LazyPM/pkg/task"
)

var taskFeedback task.ValidationFeedback

func SetTaskFeedback(feedback task.ValidationFeedback) {
	taskFeedback = feedback
}

func HandleTaskStatus(w http.ResponseWriter, r *http.Request) {
	hx := HTMX(r)

	if hx.IsHxRequest() {
		if taskFeedback.Success {
			hx.WriteString("<div>Task completed successfully!</div>")
		} else {
			hx.WriteString("<div>" + taskFeedback.Message + "</div>")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	hx.WriteJSON(taskFeedback)
}
